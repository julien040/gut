package controller

import (
	"fmt"

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
	username, err := prompt.InputLine("Username: ")
	if err != nil {
		print.Message("We can't read your input 😓", print.Error)
		return askForProfileUsername()
	}
	return username
}

// Ask a profile password to the user
func askForProfilePassword() string {
	// Ask for the password
	prompt := promptui.Prompt{
		Label: "Password",
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
		userInput, err1 := prompt.InputLine("Domain: ")
		parsedHost := getHost(userInput)
		if err1 != nil || parsedHost == "" {
			print.Message("We can't read your input 😓", print.Error)
			return askForProfileWebsite("")
		} else {
			return parsedHost
		}

	}
	return gitURL
}

func newProfile(url string) profile.Profile {
	// Create a new profile
	newProfile := profile.Profile{
		Alias:    askForProfileAlias(),
		Username: askForProfileUsername(),
		Password: askForProfilePassword(),
		Website:  askForProfileWebsite(url),
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

	fmt.Println(color.HiBlackString("ID | Alias | Username | Website"))
	for key, val := range *profiles {
		fmt.Printf(color.HiBlackString("%d. %s | %s | %s | %s"), key, val.Id, color.HiBlueString(val.Alias), color.HiBlueString(val.Username), color.HiBlueString(val.Website))
	}
}
