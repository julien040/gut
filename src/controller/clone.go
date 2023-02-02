package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/plumbing/transport"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"

	"github.com/spf13/cobra"
)

// Function executed when the user runs the clone command
func Clone(cmd *cobra.Command, args []string) {

	/* -------------------------------- Check URL ------------------------------- */
	var repo, path string
	// Check if the user provided a repo
	isUrlGiven := len(args) > 0
	// Check if the user provided a repo and a path
	isPathGiven := len(args) > 1

	if !isUrlGiven || !checkURL(args[0]) {
		// Loop until the user enters a valid URL
		str, err := prompt.InputWithValidation(
			"Which repo do you want to clone?",
			"\nUh oh, I don't think this is a valid URL 😓. Please make sure it follows the format (e.g. https://github.com/julien040/gut)",
			checkURL,
		)
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		} else {
			repo = str
		}
	} else {
		repo = args[0]
	}

	/* --------------------------------- Check path ------------------------------ */
	repoName := getRepoNameFromURL(repo)
	if isPathGiven {
		path = args[1]
	}
	path = makeValidPath(path, repoName)

	isEmpty, err := isDirectoryEmpty(path)
	if err != nil {
		exitOnError("I couldn't check if the directory is empty 😓", err)
	}
	if !isEmpty {
		// If the directory is not empty, ask the user if he wants to continue
		shouldContinue, err := prompt.InputBool("The directory is not empty. Do you want to continue? This will overwrite the existing files", true)
		if err != nil {
			exitOnError("Sorry, I can't get your answer", err)
		}
		if !shouldContinue {
			os.Exit(0)
		}
	}

	shouldConserveGitHistory, err := prompt.InputBool("Do you want to also clone the git history?", true)
	if err != nil {
		exitOnError("Sorry, I can't get your answer", err)
	}

	/* --------------------------------- Clone repo ------------------------------ */
	fmt.Printf("\nYour repo is %s and will be cloned in %s", color.GreenString(repo), color.BlueString(path))
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build a new spinner
	s.Start()
	err = executor.Clone(repo, path, shouldConserveGitHistory)
	s.Stop()
	if err != nil {
		if err.Error() == "authentication required" {
			print.Message("Oh no, this repo requires authentication 😓. Please enter your credentials", print.Info)
			cloneRepoNeedsAuth(repo, path, shouldConserveGitHistory)
		} else {
			exitOnError("Sorry but I can't clone the repo 😓", err)
		}
	} else {
		print.Message("I've successfully cloned your repo 🎉 at "+path, print.Success)
	}

}

func cloneRepoNeedsAuth(repo string, path string, shouldConserveGitHistory bool) {
	profile := selectProfile(repo, true)
	err := executor.CloneWithAuth(repo, path, profile.Username, profile.Password, shouldConserveGitHistory)
	if err == transport.ErrAuthorizationFailed {
		print.Message("Uh oh, the credentials you entered are invalid. Please try again with a different profile 😉", print.Error)
		cloneRepoNeedsAuth(repo, path, shouldConserveGitHistory)
	} else if err != nil {
		exitOnError("I can't clone the repo 😓. Please make sure you have the right permissions", err)
	} else {
		print.Message("I've successfully cloned your repo 🎉 at "+path, print.Success)
		associateProfileToPath(profile, path)
	}
}
