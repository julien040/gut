package controller

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

func Goto(cmd *cobra.Command, args []string) {
	wd := getWorkingDir()

	checkIfGitRepoInitialized(wd)

	// Add a spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	s.Prefix = "Checking if there is uncommitted changes "
	s.Start()
	// Check if there is uncommitted changes
	clean, err := executor.IsWorkTreeClean(wd)
	s.Stop()
	if err != nil {
		exitOnError("Sorry, I can't check if there is uncommitted changes", err)
	}
	if !clean {
		exitOnError("Sorry, you have uncommitted changes. Save them with \"gut save\" before going to another commit or you will lose them", nil)
	}

	// List all the commits
	commits, err := executor.ListAllCommits(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the list of commits", err)
	}

	// Check if there is commits
	if len(commits) == 0 {
		print.Message("You don't have any commit yet", print.Warning)
		return
	}

	var commit object.Commit

	// Choose the commit
	if len(args) > 0 {
		arg := args[0]
		// Check if the hash is valid
		commit, err = executor.GetCommitByHash(wd, arg)
		if err != nil {
			errMessage := err.Error()
			switch errMessage {
			case "commit not found":
				print.Message("The commit \"%s\" doesn't exist", print.Error, arg)

			case "hash must be at least 6 characters":
				print.Message("The hash \"%s\" is too short (must be at least 6 characters)", print.Error, arg)

			default:
				exitOnError("Sorry, I can't get the commit", err)

			}
			commit = chooseCommit(commits)
		}
	} else {
		commit = chooseCommit(commits)
	}

	// Get current branch for informing the user
	currentBranch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the current branch", err)
	}

	// Prompt confirmation
	res, err := prompt.InputBool("Are you sure you want to go back to \""+getTitleFromCommit(commit.Message)+"\"?", true)
	if err != nil {
		exitOnError("Sorry, I can't get your confirmation", err)
	}
	if !res {
		return
	}

	s.Prefix = "Changing your working tree to the commit " + commit.Hash.String() + " "
	s.Start()
	// Checkout the commit
	err = executor.CheckoutCommit(wd, commit.Hash.String())
	s.Stop()
	if err != nil {
		exitOnError("Sorry, I can't checkout the commit", err)
	}
	print.Message("You are now on the commit \"%s\n(%s)\"", print.Success, commit.Hash.String(), getTitleFromCommit(commit.Message))
	print.Message("To go back to the branch \"%s\", use:\n	gut switch %s", print.Optional, currentBranch, currentBranch)
	print.Message("To create a branch from this commit, use:\n	gut switch [new branch name]", print.Optional)

}
