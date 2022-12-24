/* -------------------------------------------------------------------------- */
/*     This file contains function that can be used by multiple commands.     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"fmt"
	"net/url"
	"os"

	"path/filepath"

	"github.com/julien040/gut/src/print"
)

// An helper function to exit the program if an error occurs
func exitOnError(str string, err error) {
	if err != nil {
		print.Message(str, print.Error)
		fmt.Println(err)
		os.Exit(1)
	}
}

// Return true if the string is a valid URL for git
func checkURL(str string) bool {
	parsed, err := url.Parse(str)
	if err != nil {
		return false
	} else {
		// Check if the URL has a scheme and a host
		return parsed.Scheme != "" && parsed.Host != ""
	}
}

// Return the name of the repo from the URL
func getRepoNameFromURL(str string) string {
	parsed, err := url.Parse(str)
	if err != nil {
		return ""
	} else {
		// Check if the URL has a scheme and a host
		return filepath.Base(parsed.Path)
	}
}

// Ask the user for a path to clone a repo and make sure it's valid
func makeValidPath(originalPath string, repoName string) string {
	// Check if the path exists
	path := getAbsPathFromInput(repoName, originalPath)
	if checkIfPathExist(path) {
		return path
	} else {
		// If the path doesn't exist, we ask the user if he wants to create it
		if promptUserForMakingPath(path) {
			return path
		} else {
			return makeValidPath(askForPath(repoName, "Where do you want to clone the repo?"), repoName)
		}
	}
}
