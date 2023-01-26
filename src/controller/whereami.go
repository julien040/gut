package controller

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
)

func WhereAmI(cmd *cobra.Command, args []string) {
	// Get the current directory
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory 😢", err)
	}

	// Check if the current directory is a git repository
	checkIfGitRepoInitialized(wd)

	hash, err := executor.GetHeadHash(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the current the hash of the HEAD 😢", err)
	}

	// Get the current branch
	branch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the current branch 😢", err)
	}

	fmt.Printf("HEAD is at %s on branch %s\n", color.HiGreenString(hash), color.HiGreenString(branch))

}
