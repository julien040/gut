package controller

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"time"

	"github.com/briandowns/spinner"
	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"

	"github.com/spf13/cobra"
)

func exitOnError(str string, err error) {
	if err != nil {
		/* ("For some reason, we couldn't get your input 😓", print.Error) */
		print.Message(str, print.Error)
		fmt.Println(err)
		os.Exit(1)
	}
}

// Return true if the string is a valid URL
func checkURL(str string) bool {
	parsed, err := url.Parse(str)
	if err != nil {
		return false
	} else {
		// Check if the URL has a scheme and a host
		return parsed.Scheme != "" && parsed.Host != ""
	}
}

func askForPath() string {
	res := ""
	str, err := prompt.InputLine("\nWhere do you want to clone the repo? ")
	if err != nil {
		exitOnError("For some reason, we couldn't get your input 😓", err)
	}
	if str == "." || str == "" {
		res, err = os.Getwd()

		if err != nil {
			exitOnError("For some reason, we couldn't get your input 😓", err)
		}
	} else if filepath.IsAbs(str) {
		res = str
	} else {
		// If the path is not absolute, we assume it's relative to the current directory
		res, err = filepath.Abs(str)
		if err != nil {
			exitOnError("We can't parse the path you entered. Make sure it's a valid path", err)
		}
	}
	return res
}

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
	pathExist := false
	if !isPathGiven {
		path = askForPath()

	} else {
		str, err := filepath.Abs(args[1])
		if err != nil {
			exitOnError("We can't parse the path you entered. Make sure it's a valid path", err)
		}
		path = str
	}
	for !pathExist {

		// Check if the path exists
		_, err := os.Stat(path)
		if err == nil {
			pathExist = true
		} else if errors.Is(err, os.ErrNotExist) {
			// If the path doesn't exist, we ask the user if he wants to create it

			shouldCreate, err := prompt.InputBool("The path you entered doesn't exist. Do you want to create it? \n " + color.HiBlackString(path))
			if err != nil {
				exitOnError("For some reason, we couldn't get your input 😓", err)
			}
			if shouldCreate {
				err = os.MkdirAll(path, os.ModePerm)
				if err != nil {
					exitOnError("We couldn't create the path you entered 😓. Please make sure you have the right permissions", err)
				}
				pathExist = true

			} else {
				path = askForPath()
			}
		} else if err != nil {
			exitOnError("We couldn't check if the path you entered exists 😓. Please make sure you have the right permissions", err)
		}

	}
	fmt.Printf("\nYour repo is %s and will be cloned in %s", color.GreenString(repo), color.BlueString(path))
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	err := executor.Clone(repo, path)
	s.Stop()
	if err != nil {
		exitOnError("We couldn't clone the repo 😓. Please make sure you have the right permissions", err)
	}
}
