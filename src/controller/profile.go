package controller

import (
	"fmt"

	"os"
	"path/filepath"

	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/BurntSushi/toml"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	promptui "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/profile"
	"github.com/julien040/gut/src/prompt"
	"github.com/julien040/gut/src/provider"
)

// Ask a profile alias to the user
func askForProfileAlias() string {
	alias, err := prompt.InputLine("How would you like to call this profile? ")
	if err != nil {
		print.Message("I can't read your input 😓", print.Error)
		return askForProfileAlias()
	} else if alias == "" {
		print.Message("I'm sorry but I can't create a profile with a empty name 😓", print.Error)
		return askForProfileAlias()
	} else if profile.IsAliasAlreadyUsed(alias) {
		print.Message("I'm sorry but this alias is already used 😓\nTo not confuse you, please choose another name", print.Error)
		return askForProfileAlias()
	} else {
		return alias
	}
}

// Ask a profile username to the user
func askForProfileUsername() string {
	// Ask for the username
	username, err := prompt.InputLine("What is the username? ")
	if err != nil {
		print.Message("I can't read your input 😓", print.Error)
		return askForProfileUsername()
	}
	return username
}

// Ask a profile password to the user
func askForProfilePassword() string {
	// Ask for the password
	print.Message("Your password is saved in your keychain. We can't see it 😎", print.Info)
	prompt := promptui.Prompt{
		Label: "Password or token",
		Mask:  '*',
	}
	password, err := prompt.Run()
	if err != nil {
		print.Message("Sorry, I can't read your input 😓", print.Error)
		return askForProfilePassword()
	}
	return password
}

func askForProfileWebsite(gitURL string) string {
	gitURL = getHost(gitURL)
	if gitURL == "" {
		userInput, err1 := prompt.InputLine("On which website is your repository hosted? ")
		parsedHost := getHost(userInput)
		if err1 != nil || parsedHost == "" {
			print.Message("Sorry, I can't read your input 😓", print.Error)
			return askForProfileWebsite("")
		} else if !isDomainValid(parsedHost) {
			print.Message("I think this url isn't valid 😓. Please type it again", print.Error)
			return askForProfileWebsite("")
		}
		return parsedHost

	}
	return gitURL
}

func askForProfileEmail() string {
	// Ask for the email
	email, err := prompt.InputLine("What is your email? So that I can identify your commits")
	if err != nil {
		print.Message("I can't read your input 😓", print.Error)
		return askForProfileEmail()
	} else if !isEmailValid(email) {
		print.Message("I think this email isn't valid 😓. Please type it again", print.Error)
		return askForProfileEmail()
	}
	return email
}

func newGitHubProfile() profile.Profile {
	print.Message("⚠️ None of your data will leave your computer\n", print.Info)
	// Show a spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Hand tight, I'm talking to GitHub to get a code for you..."
	s.Start()
	res, err := provider.Github_AskDeviceCode()
	s.Stop()
	if err != nil {
		exitOnError("I can't get a code from GitHub", err)
	}
	color.Black("To authenticate with GitHub, please visit the following URL and enter %s", color.WhiteString(res.UserCode))
	fmt.Println(res.VerificationURI)

	// Show a spinner
	s.Suffix = " I've done my job. Now I'm waiting for you to authenticate with GitHub... 🥱"
	s.Start()
	token, err := provider.GitHub_PollToken(res.DeviceCode)
	s.Stop()
	if err == provider.ErrExpiredToken {
		exitOnError("The code has expired. Please try again", nil)
	} else if err != nil {
		exitOnError("I can't get a token from GitHub", err)
	}
	// I need to use \n because the spinner will overwrite the message
	print.Message("I've successfully authenticated you with GitHub 🎉\n", print.Success)

	// Get the user profile
	s.Suffix = " Hand tight, Now I'm getting your profile from GitHub..."
	s.Start()
	user, err := provider.GitHub_GetUserName(token)
	s.Stop()
	if err != nil {
		exitOnError("I can't get your profile from GitHub", err)
	}
	print.Message("I've successfully got your username from GitHub 🎉\n", print.Success)

	// Get the user email
	s.Suffix = " Hand tight, Now I'm fetching your email(s) from GitHub..."
	s.Start()
	emails, err := provider.Github_GetEmails(token)
	s.Stop()
	if err != nil {
		exitOnError("I can't get your email from GitHub", err)
	}
	print.Message("I've successfully got your email(s) from GitHub 🎉\n", print.Success)
	nbrEmail := len(emails)

	var email string
	// In case no email was retrieved, we prompt the user to add one
	if nbrEmail == 0 {
		email = askForProfileEmail()
	} else if nbrEmail == 1 { // In case there is only one email, we use it
		email = emails[0]
	} else { // If several verified emails were found, we ask the user to choose one
		var emailAnswer string
		promptEmail := survey.Select{
			Message: "Which email do you want to use? (Be careful, this might be public)",
			Options: emails,
			Help:    "This email is used to sign your commits. It can be easily retrieve publicly using the git history. I recommend you to use the @users.noreply.github.com email",
		}
		err := survey.AskOne(&promptEmail, &emailAnswer)
		if err != nil {
			exitOnError("I can't read your input 😓", err)
		}
		email = emailAnswer
	}
	alias := askForProfileAlias()
	return profile.Profile{
		Alias:    alias,
		Website:  "github.com",
		Email:    email,
		Password: token,
		Username: user,
	}

}

// Create a new profile
func newProfile(url string) profile.Profile {
	// There is two type of flow: Github or anything else
	// Usign the URL, we check if the user is using Github
	// If yes, we call newGitHubProfile
	// If not, we continue the normal flow

	var newProfile profile.Profile

	// Check if the user is using Github
	host := getHost(url)

	// Case the user is using Github
	if host == "github.com" {
		newProfile = newGitHubProfile()
	} else if url != "" { // Case the user is using another platform
		newProfile = profile.Profile{
			Alias:    askForProfileAlias(),
			Username: askForProfileUsername(),
			Password: askForProfilePassword(),
			Website:  askForProfileWebsite(url),
			Email:    askForProfileEmail(),
		}
	} else { // Case we don't know which platform is used, it could be github or not so we ask the user
		var answer int

		// I use a select because BitBucket and GitLab might be added in the future
		prompt := &survey.Select{
			Message: "Which platform would you like to use ?",
			Options: []string{"GitHub", "Other"},
		}
		survey.AskOne(prompt, &answer)
		if answer == 0 { // Case the user is using Github
			newProfile = newGitHubProfile()
		} else { // Case the user is using another platform
			newProfile = profile.Profile{
				Alias:    askForProfileAlias(),
				Username: askForProfileUsername(),
				Password: askForProfilePassword(),
				Website:  askForProfileWebsite(url),
				Email:    askForProfileEmail(),
			}
		}
	}

	// Save the profile
	id := profile.AddProfile(newProfile)
	newProfile.Id = id
	print.Message("I've successfully created your profile 😎", print.Success)
	return newProfile
}

func selectProfile(gitURL string, createPossible bool) profile.Profile {
	// Get all the profiles
	profiles := profile.GetProfiles()

	// Create a slice of strings for the prompt
	var profileNames []string
	for _, val := range *profiles {
		profileNames = append(profileNames, val.Alias)
	}
	if createPossible {
		profileNames = append(profileNames, "Create a new profile")
	}
	lenProfiles := len(profileNames)
	// Ask the user to select a profile
	promptSelect := promptui.Select{
		Label: "Select a profile",
		Items: profileNames,
	}
	i, _, err := promptSelect.Run()
	if err != nil {
		exitOnError("Sorry, I can't get your answer", err)
	}
	if createPossible && i == lenProfiles-1 {
		// Create a new profile
		return newProfile(gitURL)
	} else {
		// Return the selected profile
		return (*profiles)[i]
	}

}

func Profiles(cmd *cobra.Command, args []string) {
	print.Message("Print info about the profiles \n", print.Info)
	profileSelected := selectProfile("", true)
	fmt.Printf(color.HiBlackString("Profile selected: %s \n"), color.HiBlueString(profileSelected.Alias))
	fmt.Printf(color.HiBlackString("Email: %s \n"), color.HiBlueString(profileSelected.Email))
	fmt.Printf(color.HiBlackString("Username: %s \n"), color.HiBlueString(profileSelected.Username))
	fmt.Printf(color.HiBlackString("Website: %s \n"), color.HiBlueString(profileSelected.Website))
	fmt.Printf(color.HiBlackString("Internal ID: %s \n"), color.HiBlueString(profileSelected.Id))

}

func ProfilesAdd(cmd *cobra.Command, args []string) {
	newProfile("")
}

func ProfilesList(cmd *cobra.Command, args []string) {
	// Get all the profiles
	profiles := profile.GetProfiles()

	if len(*profiles) == 0 {
		print.Message("You don't have any profile yet 😓 \nCreate one with gut profile add", print.Info)
		return
	} else {

		fmt.Println(color.HiBlackString("ID | Alias | Username | Website | Email"))
		for key, val := range *profiles {
			fmt.Printf(color.HiBlackString("%d. %s | %s | %s | %s | %s \n"), key, val.Id, color.HiBlueString(val.Alias), color.HiBlueString(val.Username), color.HiBlueString(val.Website), color.HiBlueString(val.Email))
		}
	}
}
func associateProfileToPath(profileArg profile.Profile, path string) {
	// Set Profile Data in git config
	executor.SetUserConfig(path, profileArg.Username, profileArg.Email)

	// Get current date
	currentDate := time.Now().Format("2006-01-02 15:04:05")

	// Check if file exists
	pathToWrite := filepath.Join(path, ".gut")
	if _, err := os.Stat(pathToWrite); os.IsNotExist(err) {
		// Create file
		f, err := os.Create(pathToWrite)
		if err != nil {
			exitOnError("I can't create the file .gut at "+pathToWrite, err)
		}
		f.Close()
	}

	// Open file in write mode
	f, err := os.OpenFile(pathToWrite,
		os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		exitOnError("I can't open the file .gut at "+pathToWrite, err)
	}

	// Close file at the end of the function
	defer f.Close()

	// Create the schema
	profileIDSchema := profile.SchemaGutConf{
		ProfileID: profileArg.Id,
		UpdatedAt: currentDate,
	}

	// Encode ID in TOML and write it in .gut file
	t := toml.NewEncoder(f)
	err = t.Encode(profileIDSchema)
	if err != nil {
		exitOnError("I can't write the profile ID in the file .gut at "+pathToWrite, err)
	}

	// Change Git Config
	executor.SetUserConfig(path, profileArg.Username, profileArg.Email)

}
