package controller

import (
	"os"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func Switch(cmd *cobra.Command, args []string) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't find the current working directory", err)
	}
	// Check if the current directory is a git repository
	checkIfGitRepoInitialized(wd)

	// Input the branch name
	var branchName string
	if len(args) == 0 {
		res, err := prompt.InputLine("Switch to branch: ")
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		}
		branchName = res
	} else {
		branchName = args[0]
	}

	// Check if the branch exists
	exists, err := executor.CheckIfBranchExists(wd, branchName)
	if err != nil {
		exitOnError("I can't check if the branch exists", err)
	}

	// If the branch doesn't exist, ask the user if he wants to create it
	if !exists {
		res, err := prompt.InputBool("Uh oh, the branch doesn't exist. Do you want me to create it?", true)
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		}
		if res {
			err = executor.CreateBranch(wd, branchName)
			if err != nil {
				exitOnError("My bad, I can't create the branch", err)
			}
		} else {
			print.Message("Okay, I won't create the branch", print.Info)
			return
		}
	} else { // If the branch exists, switch to it

		// Check if the working tree is clean
		clean, err := executor.IsWorkTreeClean(wd)
		if err != nil {
			exitOnError("I can't check if there are uncommitted changes", err)
		}
		// If not clean, ask the user if he wants to continue because the changes might be lost
		if !clean {
			res, err := prompt.InputBool("Uh oh, there are uncommitted changes. They might be lost if you switch branches. Do you want to continue?", false)
			if err != nil {
				exitOnKnownError(errorReadInput, err)
			}
			if !res {
				print.Message("Okay, I won't switch branches", print.Info)
				return
			}
		}
		err = executor.CheckoutBranch(wd, branchName)
		if err != nil {
			exitOnError("I can't switch to the branch "+branchName, err)
		}

	}
	print.Message(`I switched to the branch "`+branchName+`" successfully 🎉`, print.Success)

}
