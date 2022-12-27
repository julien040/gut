package controller

import (
	"os"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/spf13/cobra"
)

func Init(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Oups, something went wrong while getting the current working directory", err)
	}
	if executor.IsPathGitRepo(wd) {
		exitOnError("Oups, this directory is already a git repository. Delete the .git folder if you want to initialize a new repository", nil)
	}
	err = executor.Init(wd)
	if err != nil {
		exitOnError("Oups, something went wrong while initializing the repository", err)
	}
	profile := selectProfile("", true)
	associateProfileToPath(profile, wd)
	_, err = executor.Commit(wd, "🎉 Initial commit from Gut")
	if err != nil {
		exitOnError("Oups, something went wrong while creating the first commit", err)
	}
	print.Message("Yeah, your repository is ready to go!", print.Success)

}
