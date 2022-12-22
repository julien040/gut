package controller

import (
	"fmt"

	"github.com/fatih/color"

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
			"\nUh oh, we can't parse the repo you entered. Make sure it's a valid URL (e.g. https://example.com/repo.git) ",
			checkURL,
		)
		if err != nil {
			exitOnError("For some reason, we couldn't get your input 😓", err)
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

	/* --------------------------------- Clone repo ------------------------------ */
	fmt.Printf("\nYour repo is %s and will be cloned in %s", color.GreenString(repo), color.BlueString(path))
	/* s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build a new spinner
	s.Start()
	err := executor.Clone(repo, path)
	s.Stop()
	if err == transport.ErrAuthenticationRequired {
		print.Message("We don't have the right credentials to clone the repo 😓. Please make sure you have the right permissions", print.Error)
		os.Exit(1)
	} else if err != nil {
		exitOnError("We couldn't clone the repo 😓. Please make sure you have the right permissions", err)
	} */

	// Todo: Add auth handler for private repos
	// Todo: Add degit support
}
