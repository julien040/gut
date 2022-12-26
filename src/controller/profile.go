package controller

import (
	"fmt"

	"os"
	"path/filepath"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/profile"
	"github.com/julien040/gut/src/prompt"
	promptui "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Ask a profile alias to the user
func askForProfileAlias() string {
	alias, err := prompt.InputLine("Alias: ")
	if err != nil {
		print.Message("We can't read your input 😓", print.Error)
		return askForProfileAlias()
	} else if alias == "" {
		print.Message("The alias can't be empty 😓", print.Error)
		return askForProfileAlias()
	} else {
		return alias
	}
}

// Ask a profile username to the user
func askForProfileUsername() string {
	// Ask for the username
	username, err := prompt.InputLine("Login username: ")
	if err != nil {
		print.Message("We can't read your input 😓", print.Error)
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
		print.Message("We can't read your input 😓", print.Error)
		return askForProfilePassword()
	}
	return password
}

func askForProfileWebsite(gitURL string) string {
	gitURL = getHost(gitURL)
	if gitURL == "" {
		userInput, err1 := prompt.InputLine("Website: ")
		parsedHost := getHost(userInput)
		if err1 != nil || parsedHost == "" {
			print.Message("We can't read your input 😓", print.Error)
			return askForProfileWebsite("")
		} else if !isDomainValid(parsedHost) {
			print.Message("We think this URL is not valid 😓. Please type it again", print.Error)
			return askForProfileWebsite("")
		}
		return parsedHost

	}
	return gitURL
}

func askForProfileEmail() string {
	// Ask for the email
	email, err := prompt.InputLine("Email (for commit): ")
	if err != nil {
		print.Message("We can't read your input 😓", print.Error)
		return askForProfileEmail()
	} else if !isEmailValid(email) {
		print.Message("We think this email is not valid 😓. Please type it again", print.Error)
		return askForProfileEmail()
	}
	fmt.Println(email)
	return email
}

func newProfile(url string) profile.Profile {
	// Create a new profile
	newProfile := profile.Profile{
		Alias:    askForProfileAlias(),
		Username: askForProfileUsername(),
		Password: askForProfilePassword(),
		Website:  askForProfileWebsite(url),
		Email:    askForProfileEmail(),
	}
	// Save the profile
	id := profile.AddProfile(newProfile)
	newProfile.Id = id
	print.Message("Profile created successfully 🎉 \n", print.Success)
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
		exitOnError("We can't read your input 😓", err)
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
func associateProfileToPath(profileID string, path string) {
	// Get current date
	currentDate := time.Now().Format("2006-01-02 15:04:05")
	// Check if file exists
	pathToWrite := filepath.Join(path, ".gut")
	if _, err := os.Stat(pathToWrite); os.IsNotExist(err) {
		// Create file
		f, err := os.Create(pathToWrite)
		if err != nil {
			exitOnError("We can't create the file .gut at "+pathToWrite, err)
		}
		f.Close()
	}
	// Open file in write mode
	f, err := os.OpenFile(pathToWrite,
		os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		exitOnError("We can't open the file .gut at "+pathToWrite, err)
	}
	// Close file at the end of the function
	defer f.Close()

	// Create the schema
	profileIDSchema := SchemaGutConf{
		ProfileID: profileID,
		UpdatedAt: currentDate,
	}

	// Encode ID in TOML
	t := toml.NewEncoder(f)
	err = t.Encode(profileIDSchema)
	if err != nil {
		exitOnError("We can't encode in TOML", err)
	}

	// Write profile id in the file
	/* _, err = f.WriteString(profileID)
	if err != nil {
		exitOnError("We can't write in the file .gutconf at "+pathToWrite, err)
	} */

}

type SchemaGutConf struct {
	ProfileID string `toml:"profile_id"`
	UpdatedAt string `toml:"updated_at"`
}
