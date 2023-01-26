/* -------------------------------------------------------------------------- */
/*     This file contains function that can be used by multiple commands.     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
		print.Message("If this error persists, please open an issue on GitHub: https://github.com/julien040/gut/issues/new", print.None)
	}

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

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
//
// Open the default browser to the specified URL.
func openInBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err

}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		exitOnError("Oups, something went wrong while I was getting the current directory", err)
	}
	return dir
}

// Check if Repo is in detached HEAD state
//
// If it is, ask the user if he wants to go back to a branch or create a new one from the commit
func checkIfDetachedHead(wd string) {
	detached, err := executor.IsDetachedHead(wd)
	if err != nil {
		exitOnError("Oups, something went wrong while I was checking if the repo is in detached HEAD state", err)
	}
	if !detached {
		return
	}
	print.Message("Oups, it seems like you are in lookout mode (detached HEAD state)", print.None)

	// Ask the user if he wants to checkout a branch or create a new one
	options := []string{"1) Switch back to a branch", "2) Create a new branch from this commit"}
	val, err := prompt.InputSelect("What do you want to do?", options)
	if err != nil {
		exitOnError("Oups, something went wrong while I was getting the input", err)
	}
	switch val {
	case "1) Switch back to a branch":
		// Choose a branch to checkout
		branches, err := executor.ListBranches(wd)
		if err != nil {
			exitOnError("Oups, something went wrong while I was listing the branches", err)
		}
		branch, err := prompt.InputSelect("Which branch do you want to checkout?", branches)
		if err != nil {
			exitOnError("Oups, something went wrong while I was getting the input", err)
		}
		// Checkout the branch
		err = executor.CheckoutBranch(wd, branch)
		if err != nil {
			exitOnError("Oups, something went wrong while I was checking out the branch", err)
		}
		print.Message("You are now on branch %s and you can continue working 🎉", print.Success, branch)
		os.Exit(0)

	case "2) Create a new branch from this commit":
		// Ask the user for the name of the new branch
		var promptUser func() string
		promptUser = func() string {
			branchName, err := prompt.InputLine("How do you want to name the new branch?")
			if err != nil {
				exitOnError("Oups, something went wrong while I was getting the input", err)
			}
			if branchName == "" {
				print.Message("Oups, you need to enter a name for the branch", print.Error)
				return promptUser()
			}
			// Check if the branch already exists
			exists, err := executor.CheckIfBranchExists(wd, branchName)
			if err != nil {
				exitOnError("Oups, something went wrong while I was checking if the branch already exists", err)
			}
			if exists {
				print.Message("Oups, a branch with this name already exists", print.Error)
				return promptUser()
			}
			return branchName
		}
		branchName := promptUser()

		// Create the new branch
		err = executor.CreateBranch(wd, branchName)
		if err != nil {
			exitOnError("Oups, something went wrong while I was creating the new branch", err)
		}
		print.Message("You are now on branch %s and you can continue working 🎉", print.Success, branchName)
		os.Exit(0)

	default:
		exitOnError("Oups, something went wrong while I was getting the input", nil)

	}

}
