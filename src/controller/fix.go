package controller

import (
	"os"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"

	"github.com/spf13/cobra"
)

func Fix(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't find the current working directory", err)
	}

	options := []string{
		"1) I've made a typo in the last commit message",
		"2) I've committed on the wrong branch",
		"3) I want to discard changes since last commit",
		"4) I want to go back to a previous commit",
		"5) I forgot to add a change in the last commit",
		"Cancel",
	}

	checkIfGitRepoInitialized(wd)

	print.Message("You made a mess, huh? Let's fix it!", print.None)
	res, err := prompt.InputSelect("What do you want to fix?", options)
	if err != nil {
		exitOnError("Sorry, I can't read your answer", err)
	}

	switch res {
	case "1) I've made a typo in the last commit message":
		amendMessage(wd)
	case "2) I've committed on the wrong branch":
		cherryPick()

	case "3) I want to discard changes since last commit":
		print.Message("I have a command for that: gut undo", print.Info)

	case "4) I want to go back to a previous commit":
		print.Message("I have a command for that: gut revert", print.Info)

	case "5) I forgot to add a change in the last commit":
		amendCommit()

	case "Cancel":

	default:
		print.Message("Sorry, I don't know how to fix that", print.Error)

	}

}

/*
Try to amend the last commit message
It will only work if the commit hasn't been pushed yet
*/
func amendMessage(path string) {
	// Check if Git is installed
	installed := executor.IsGitInstalled()
	if !installed {
		exitOnError("Sorry, I can't find Git on your computer", nil)
	}

	/* // Get remotes
	remotes, err := executor.ListRemote()
	if err != nil {
		exitOnError("Sorry, I can't get the remotes", err)
	}
	if len(remotes) != 1 {
		exitOnError("Sorry, I can't handle multiple remotes", nil)
	}

	remote := remotes[0] */

	// Get head commit
	head, err := executor.GetHeadHash(path)
	if err != nil {
		exitOnError("Sorry, I can't get the head commit", err)
	}

	// Check if there are uncommitted changes
	clean, err := executor.IsWorkTreeClean(path)
	if err != nil {
		exitOnError("Sorry, I can't check if there are uncommitted changes", err)
	}
	if !clean {
		res, err := prompt.InputBool("You have uncommitted changes. If you continue, I will include them in the last commit. Are you sure you want to continue?", true)
		if err != nil {
			exitOnError("Sorry, I can't read your answer", err)
		}
		if !res {
			return
		}
	}

	// Check if commit has been pushed
	contains := executor.GitRemoteContainsHash(head)
	if contains {
		exitOnError("Sorry, I can't change the last commit message because it has been sync with the remote\nIt might break other people's work", nil)

	}
	// Get the new commit message

	print.Message("\nLet's write the new commit message", print.None)

	message := promptCommitMessage("", "")

	// Prompt a confirmation
	res, err := prompt.InputBool("Are you sure you want me to change the last commit message?", true)
	if err != nil {
		exitOnError("Sorry, I can't read your answer", err)
	}
	if !res {
		return
	}

	// Amend the commit
	err = executor.GitCommitAmend(message)
	if err != nil {
		exitOnError("Sorry, I can't amend the last commit", err)
	}
	print.Message("I've successfully changed the last commit message", print.Success)

}

func cherryPick() {
	print.Message("For now, I can't move a commit to another branch", print.Error)
	print.Message("Aha", print.None)
}

func amendCommit() {
	print.Message("For now, I can't add a change to the last commit", print.Error)
	print.Message("Aha", print.None)
}
