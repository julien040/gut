package controller

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

// Returns the index of the latest commit that has been pushed
// Returns -1 if no commit has been pushed
//
// For example, if we have the following commits:
//   - a5be1b2 (not pushed)
//   - 3b4c5d6 (not pushed)
//   - 7e8f9g0 (not pushed)
//   - 1a2b3c4 (pushed)
//   - 4d5e6f7 (pushed)
//
// The function will return 3
func getIndexLatestCommitPushed(commits []object.Commit) int {
	for i, commit := range commits {
		if executor.GitRemoteContainsHash(commit.Hash.String()) {
			return i
		}
	}
	return -1
}

// Squash squashes all commits to a specific commit
// It will ask the user to select a commit
// Then it will squash all commits to this commit
//
// To do so, it will:
//   - Reset softly to the selected commit in order to leave the changes in the working directory but move the HEAD
//   - Amend the commit with the new message and the new content

func Squash(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory", err)
	}

	// Check if the current directory is a git repository
	checkIfGitRepoInitialized(wd)

	// Check if the HEAD is detached
	checkIfDetachedHead(wd)

	// Check if the user has a valid configuration
	verifUserConfig(wd)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Listing all commits... "
	s.Start()
	// List all commits from the current branch
	commits, err := executor.ListCommit(wd)
	s.Stop()
	if err != nil {
		exitOnError("Sorry, I can't list the commits", err)
	}

	// Get the index of the latest commit that has been pushed
	// Because gut doesn't allow rewriting the history, we can't squash a commit that has been pushed
	s.Prefix = "Checking if your commits have been pushed... "
	s.Start()
	indexLatestCommitPushed := getIndexLatestCommitPushed(commits)
	s.Stop()

	// It means the user has already pushed all his commits
	// We can't squash them because it would rewrite the history
	// So we exit the program
	if indexLatestCommitPushed == 0 {
		print.Message("All your commits have been pushed. Because I won't allow any rewrite of the history, I can't squash your commits", print.Warning)
		os.Exit(0)
	}

	// Check if there is enough commits to squash
	if len(commits) < 2 || indexLatestCommitPushed == 1 {
		print.Message("You don't have enough commits to squash. Please make at least 2 commits before squashing them", print.Warning)
		os.Exit(0)
	}

	// The object.Commit to squash
	// Gut will soft reset to this commit
	// And amend it with a new commit and a new message
	var commitToSquash object.Commit

	// Prompt the user to choose a commit to squash
	promptCommitToSquash := func() object.Commit {
		// We start at 1 because we don't want to squash the latest commit
		return chooseCommit(commits[1:indexLatestCommitPushed])
	}

	// If the user has passed a commit hash as argument
	if len(args) > 0 {
		// The commit hash passed as argument
		commitArg := args[0]

		// We retrieve the commit object from the short hash
		commitToSquash, err = executor.GetCommitByHash(wd, commitArg)

		// Case 1: The commit doesn't exist
		if err != nil {
			print.Message("I can't find the commit %s. Please choose a commit from the list below", print.Warning, commitArg)
			commitToSquash = promptCommitToSquash()

		} else if executor.GitRemoteContainsHash(commitToSquash.Hash.String()) { // Case 2: The commit has already been pushed
			print.Message("The commit %s has already been pushed. I don't allow rewriting of the history. Please choose a commit from the list below", print.Warning, commitArg)
			commitToSquash = promptCommitToSquash()
		}

	} else { // If the user hasn't passed a commit hash as argument
		commitToSquash = promptCommitToSquash()
	}

	print.Message("Choose a new message for the commit", print.Info)
	// Choose a new message for the commit
	newMessage := promptCommitMessage("", "")

	res, err := prompt.InputBool("Are you sure you want to squash all commits to "+commitToSquash.Hash.String()+"?", false)
	if err != nil {
		exitOnError("Sorry, I can't prompt the user", err)
	}
	if !res {
		os.Exit(0)
	}

	// Soft reset to the commit to squash
	s.Prefix = "Soft resetting to the commit to squash... "
	s.Start()
	err = executor.SoftReset(wd, commitToSquash)
	s.Stop()
	if err != nil {
		exitOnError("Sorry, I can't soft reset to the commit to squash", err)
	}

	// Amend the commit with the new message
	s.Prefix = "Amending the commit... "
	s.Start()
	err = executor.GitCommitAmend(newMessage)
	s.Stop()
	if err != nil {
		exitOnError("Sorry, I can't amend the commit", err)
	}

	print.Message("The commit has been successfully squashed", print.Success)

}
