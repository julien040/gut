package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/AlecAivazis/survey/v2"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func getTitleFromCommit(message string) string {
	// Split the message by new line
	lines := strings.Split(message, "\n")
	// Get the first line
	return lines[0]

}

func chooseCommit(commits []object.Commit) object.Commit {
	// Create the list of commits
	choices := []string{}
	for _, commit := range commits {
		choices = append(choices, fmt.Sprintf("%s created by %s on %s" /* color.HiYellowString(commit.Hash.String()), */, color.HiCyanString(getTitleFromCommit(commit.Message)), commit.Author.Name, commit.Author.When.Format("Mon Jan 2 15:04:05")))
	}
	// Ask the user to choose a commit
	qs := &survey.Select{

		Message: "Choose a commit to undo",
		Options: choices,
	}
	// Get the answer
	var answer int
	err := survey.AskOne(qs, &answer)
	if err != nil {
		exitOnError("Sorry, I can't get your answer 😢", err)
	}

	return commits[answer]
}

func Revert(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory 😢", err)
	}

	// Check if the current directory is a git repository
	checkIfGitRepoInitialized(wd)

	// Check if Git CLI is installed
	checkIfGitInstalled()

	print.Message("Undo reverts your working tree to a commit of your choice", print.Info)

	// Check if the working tree is clean (no uncommitted changes)
	wtClean, err := executor.IsWorkTreeClean(wd)
	if err != nil {
		exitOnError("Sorry, I can't check if the working tree is clean 😢", err)
	}
	if !wtClean {
		shouldContinue, err := prompt.InputBool("Your working tree is not clean. Changes will be lost. Do you want to continue?", false)
		if err != nil {
			exitOnError("Sorry, I can't get your answer 😢", err)
		}
		if !shouldContinue {
			print.Message("Bye 👋", print.Success)
			return
		}
	}
	// Use go git to get the list of commits sorted by date (newest first)
	commits, err := executor.ListCommit(wd)
	if err != nil {
		exitOnError("Sorry, I can't list the commits 😢", err)
	}
	// Case if there is only one commit or no commit
	if len(commits) < 2 {
		exitOnError("Sorry, there is no commit to undo 😢", nil)
	}

	// Prompt the user to choose a commit
	commit := chooseCommit(commits)
	fmt.Printf("I will revert the commit to %s created by %s on %s \n\n", color.HiCyanString(getTitleFromCommit(commit.Message)), color.HiCyanString(commit.Author.Name), commit.Author.When.Format("Mon Jan 2 15:04:05 2006"))

	err = executor.GitRevert(commit.Hash.String())
	if err != nil {
		exitOnError("Sorry, I can't revert the commit. An error occured while calling 'git revert --no-edit "+commit.Hash.String()+"' 😢", err)
	}
	err = executor.AddAll(wd)
	if err != nil {
		exitOnError("Sorry, I can't add all the files 😢", err)
	}

	print.Message("I've successfully to "+getTitleFromCommit(commit.Message), print.Success)

}
