package controller

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
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
	var refArg string
	if len(args) == 0 {
		res, err := prompt.InputLine("Switch to: ")
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		}
		refArg = res
	} else {
		refArg = args[0]
	}

	// Define a reusable spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	// Check if the branch exists
	exists, err := executor.CheckIfBranchExists(wd, refArg)
	if err != nil {
		exitOnError("I can't check if the branch exists", err)
	}

	// If the branch doesn't exist, we check if the user wants to switch to a commit
	if !exists {
		// Check if the user wants to switch to a commit
		commit, err := executor.GetCommitByHash(wd, refArg)
		// If the commit exists, we switch to it
		if err == nil {

			// Check if the working tree is clean
			s.Prefix = "Checking if the working tree is clean... "
			s.Start()
			clean, err := executor.IsWorkTreeClean(wd)
			s.Stop()
			if err != nil {
				exitOnError("I can't check if there are uncommitted changes", err)
			}
			// If not clean, we return an error
			if !clean {
				exitOnError("Sorry, you have uncommitted changes. Save them with \"gut save\" before going to another commit or you will lose them", nil)
			}

			// Switch to the commit
			s.Prefix = "Switching to the commit... "
			s.Start()
			err = executor.CheckoutCommit(wd, commit.Hash.String())
			s.Stop()
			if err != nil {
				exitOnError("I can't switch to the commit", err)
			}
			print.Message("I've successfully switched to the commit %s", print.Info, commit.Hash.String())
			return
		}

		res, err := prompt.InputBool("Uh oh, the branch doesn't exist. Do you want me to create it?", true)
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		}
		if res {
			err = executor.CreateBranch(wd, refArg)
			if err != nil {
				exitOnError("My bad, I can't create the branch", err)
			}
		} else {
			print.Message("Okay, I won't create the branch", print.Info)
			return
		}
	} else { // If the branch exists, switch to it

		// Check if the branch is the current branch
		currentBranch, err := executor.GetCurrentBranch(wd)
		if err != nil {
			exitOnError("I can't get the current branch", err)
		}
		if currentBranch == refArg {
			print.Message("You are already on the branch "+refArg, print.Success)
			return
		}

		// Check if the working tree is clean
		s.Prefix = "Checking if the working tree is clean... "
		s.Start()

		// Set to true when a stash is created
		// We need to pop the stash at the end of the function
		mustStashPop := false

		// Check if the working tree is clean
		clean, err := executor.IsWorkTreeClean(wd)
		s.Stop()
		if err != nil {
			exitOnError("I can't check if there are uncommitted changes", err)
		}
		// If not clean, ask the user if he wants to continue because the changes might be lost
		if !clean {
			/*
				Because the working tree is not clean, we have two options:
					- stash the changes
					- discard the changes

				We will ask the user which option he wants to use
			*/
			print.Message("Uh oh, there are uncommitted changes", print.Warning)

			res, err := prompt.InputSelect("What do you want to do?", []string{"Keep the changes", "Discard the changes"})
			if err != nil {
				exitOnKnownError(errorReadInput, err)
			}

			// If he wants to keep the changes, we stash them
			// Also, we set a flag so that we can pop the stash at the end of the function
			if res == "Keep the changes" {
				executor.GitStash()
				mustStashPop = true
			}
		}
		s.Prefix = "Switching to the branch " + refArg + " "
		s.Start()

		err = executor.CheckoutBranch(wd, refArg)

		// If we stashed the changes, we need to pop the stash
		if mustStashPop {
			executor.GitStashPop()
		}
		s.Stop()
		if err != nil {
			exitOnError("I can't switch to the branch "+refArg, err)
		}

	}
	print.Message(`I switched to the branch "`+refArg+`" successfully 🎉`, print.Success)

}
