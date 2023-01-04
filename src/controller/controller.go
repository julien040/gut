/* -------------------------------------------------------------------------- */
/*     This file contains function that can be used by multiple commands.     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

// An helper function to exit the program if an error occurs
func exitOnError(str string, err error) {
	fmt.Println("")
	if str != "" {
		print.Message(str, print.Error)
	}
	if err != nil {
		fmt.Printf("Error message: %s\n", err)
	}
	print.Message("If this error persists, please open an issue on GitHub: https://github.com/julien040/gut/issues/new", print.None)
	os.Exit(1)
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
		base := filepath.Base(parsed.Path)
		// Remove the extension
		return strings.TrimSuffix(base, filepath.Ext(base))
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

func getHost(str string) string {
	parsed, err := url.Parse(str)
	if err != nil {
		return ""
	} else {
		// If the URL has no scheme, Parse won't find the host
		if parsed.Scheme == "" {
			return strings.Split(str, "/")[0]
		}
		return parsed.Host
	}
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil

}

func isDomainValid(domain string) bool {
	// Regex from https://stackoverflow.com/questions/10306690/what-is-a-regular-expression-which-will-match-a-valid-domain-name-without-a-subd
	// I have no idea how it works
	reg := regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
	return reg.MatchString(domain)

}

// Check if the path is a git repo
//
// If it's not, ask the user if he wants to initialize a new repo
// If he doesn't, exit the program. If he does, initialize a new repo
func checkIfGitRepoInitialized(path string) {
	// Check if the path is a git repo
	if !executor.IsPathGitRepo(path) {
		print.Message("Oups, this directory is not a git repository", print.Info)
		val, err := prompt.InputBool("Do you want to initialize a new repository?", true)
		if err != nil {
			exitOnError("Oups, something went wrong while getting the input", err)
		}
		if val {
			Init(&cobra.Command{}, []string{})
		} else {
			os.Exit(0)
		}
	}
}

func checkIfGitInstalled() {
	installed := executor.IsGitInstalled()
	if !installed {
		exitOnError("Oups, it seems like git is not installed on your computer and I need it to run this command.\nPlease follow the instructions on https://git-scm.com/book/en/v2/Getting-Started-Installing-Git to install it", nil)
	}
}
