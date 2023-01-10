package controller

import (
	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"

	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func History(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't the current working directory", err)
	}
	checkIfGitRepoInitialized(wd)

	commits, err := executor.ListCommit(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the list of commits", err)
	}
	if len(commits) == 0 {
		print.Message("You don't have any commit yet", print.Warning)
		return
	}
	commit := chooseCommit(commits)
	title := getTitleFromCommit(commit.Message)

	color.Black("\n\nCommit \"%s\" \nmade by %s on %s", color.GreenString(title), color.GreenString(commit.Author.Name), color.GreenString(commit.Author.When.Format("Mon Jan 2 2006 15:04:05 ")))

	color.Black("\nHash: %s", commit.Hash.String())

	color.Black("Message: \n%s", color.WhiteString(commit.Message))

}
