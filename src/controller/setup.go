package controller

import (
	"os"

	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func SetupAuth(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Oups, something went wrong while getting the current working directory", err)
	}
	checkIfGitRepoInitialized(wd)
	id := getProfileIDFromPath(wd)
	if id != "" {
		val, err := prompt.InputBool("This repository is already associated with a profile. Do you want to change it?", true)
		if err != nil {
			exitOnError("Oups, something went wrong while getting the input", err)
		}
		if !val {
			os.Exit(0)
		}

	}
	profile := selectProfile("", true)
	associateProfileToPath(profile, wd)
	print.Message("Whoo, your repository is now associated with the profile "+profile.Alias+" 🎉", print.Success)
}
