package controller

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
)

/*
Show changes between commits, and the working tree (requires git)
If no arguments are given, the diff will be between the current commit and the working tree.
If one argument is given, the diff will be between the given commit and the working tree.
If two arguments are given, the diff will be between the two given commits., */

func Diff(cmd *cobra.Command, args []string) {
	wd := getWorkingDir()

	checkIfGitRepoInitialized(wd)

	// Check if git is installed
	checkIfGitInstalled()

	switch len(args) {
	case 0:
		diffCurrentCommit(wd)
	case 1:
		diffCommitAndWorkingTree(wd, args[0])
	case 2:
		diffTwoCommits(wd, args[0], args[1])
	default:
		print.Message("Sorry, I don't know how to diff %d commits", print.Error, len(args))

	}

}

func diffCurrentCommit(wd string) {
	// Add all the files in the working tree to the index
	err := executor.AddAll(wd)
	if err != nil {
		exitOnError("Sorry, I can't add all the files to the index", err)
	}
	empty, _ := executor.GitDiffIndexHead()
	if empty {
		print.Message("No changes to commit", print.Warning)
	}

}

func diffCommitAndWorkingTree(wd string, commit string) {
	err := executor.AddAll(wd)
	if err != nil {
		exitOnError("Sorry, I can't add all the files to the index", err)
	}
	empty, _ := executor.GitDiffIndexCommit(commit)
	if empty {
		print.Message("No changes on the working tree from %s", print.Warning, commit)
	}
}

func diffTwoCommits(wd string, commit1 string, commit2 string) {
	// This command is slightly different. If using GitHub or GitLab, it will show the difference with the browser.
	// If using a local repository, it will show the difference with the terminal.

	// Check if the two commits have been pushed
	existsOnRemote1 := executor.GitRemoteContainsHash(commit1)
	existsOnRemote2 := executor.GitRemoteContainsHash(commit2)

	// Check if the two commits are branches
	isBranch1, _ := executor.CheckIfBranchExists(wd, commit1)
	if isBranch1 {
		existsOnRemote1 = true
	}
	isBranch2, _ := executor.CheckIfBranchExists(wd, commit2)
	if isBranch2 {
		existsOnRemote2 = true
	}

	// Check if the two commits are HEAD
	if commit1 == "HEAD" {
		existsOnRemote1 = true
	}
	if commit2 == "HEAD" {
		existsOnRemote2 = true
	}

	// If both refs are pushed, we can use the browser
	if existsOnRemote1 && existsOnRemote2 {

		// Get the host
		platform, remoteURL, err := getPlatformUsed(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the platform used", err)
		}
		switch platform {
		case "github.com":
			// https://docs.github.com/en/pull-requests/committing-changes-to-your-project/viewing-and-comparing-commits/comparing-commits
			urlToOpen := fmt.Sprintf("%s/compare/%s..%s", remoteURL, commit1, commit2)
			print.Message("Opening %s", print.Optional, color.BlueString(urlToOpen))
			openInBrowser(urlToOpen)

		case "gitlab.com":
			// https://stackoverflow.com/a/50070145/15573415
			urlToOpen := fmt.Sprintf("%s/compare?from=%s&to=%s", remoteURL, commit1, commit2)
			print.Message("Opening %s", print.Optional, color.BlueString(urlToOpen))
			openInBrowser(urlToOpen)

		default:
			empty, _ := executor.GitDiffRef(commit1, commit2)
			if empty {
				print.Message("No changes between %s and %s", print.Warning, commit1, commit2)
			}
		}

	} else {
		// If one of the commits is not pushed, we can't use the browser
		empty, _ := executor.GitDiffRef(commit1, commit2)
		if empty {
			print.Message("No changes between %s and %s", print.Warning, commit1, commit2)
		}

	}
}
