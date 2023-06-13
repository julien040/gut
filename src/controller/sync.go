package controller

import (
	"fmt"
	"os"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/profile"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"

	"github.com/briandowns/spinner"
)

func Sync(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory 😢", err)
	}
	// Check if the repository is initialized
	checkIfGitRepoInitialized(wd)

	// Check if the head is detached
	checkIfDetachedHead(wd)

	// Check if there is uncommited changes
	// If there is, we ask the user to commit them
	// If there is not, we exit

	/* clean, err := executor.IsWorkTreeClean(wd)
	if err != nil {
		exitOnError("Sorry, I can't check if there are uncommited changes 😢", err)
	}
	if !clean {
		exitOnError("Uh oh, there are uncommited changes. Please commit them before syncing with `gut save`", nil)
	} */

	// Get the remote
	// If there is no remote, we ask the user to add one
	// If there is more than one remote, we ask the user to select one
	// If there is only one remote, we use it
	remote, err := getRemote(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the remote 😢", err)
	}

	fmt.Println(remote)

	behind, err := executor.GitBehindCommitsCount()
	if err != nil {
		exitOnError("I can't get the number of commit behind", err)
	}

	ahead, err := executor.GitAheadCommitsCount()
	if err != nil {
		exitOnError("I can't get the number of commit ahead", err)
	}

	print.Message("Ahead %d | Behind %d", print.Info, ahead, behind)

	// We extract the code to use it recursively in case of error
	/* syncRepo(wd, remote, false) */
}

func syncRepo(path string, remote executor.Remote, requestProfile bool) error {
	/*
		This must get refactored for better readability
	*/
	var profileLocal profile.Profile

	// If it's not the first time syncRepo is called, we ask the user to select a profile
	if requestProfile {
		profileLocal = selectProfile("", true)
	} else { // Else, we get the profile from the path
		profilePath, err := profile.GetProfileFromPath(path)
		if err != nil { // If there is no profile associated with the path, we ask the user to select one
			profileLocal = selectProfile("", true)
		} else { // Else, we use the profile associated with the path
			profileLocal = profilePath
		}
	}
	// Once we have the profile, we pull the repository
	spinnerPull := spinner.New(spinner.CharSets[9], 100)
	spinnerPull.Prefix = "Pulling the repository from " + remote.Name + " "
	spinnerPull.Start()
	err := executor.Pull(path, remote.Name, profileLocal.Username, profileLocal.Password)
	spinnerPull.Stop()
	// If the credentials are wrong, we ask the user to select another profile
	if err == transport.ErrAuthorizationFailed || err == transport.ErrAuthenticationRequired {
		print.Message("Uh oh, your credentials are wrong 😢. Please select another profile.", print.Error)
		return syncRepo(path, remote, true)

	} else if err == git.ErrNonFastForwardUpdate {
		// Try to do Git Pull if installed

		// Check if git is installed
		gitInstalled := executor.IsGitInstalled()
		if gitInstalled {
			print.Message("I couldn't resolve this as a fast-forward merge. I'll try to do a git pull instead with the git cli", print.None)

			// Set git config to rebase if the user didn't set anything else
			err := executor.GitConfigPullRebaseTrueOnlyIfNotSet()
			if err != nil {
				exitOnError("Sorry, I can't set the git config to rebase the merge", err)
			}

			err = executor.GitPull(profileLocal.Username, profileLocal.Password, remote.Url)
			if err == nil {
				print.Message("Pull successful 🎉", print.Success)
				// We then push
				shouldReturn, returnValue := push(remote, path, profileLocal)
				if shouldReturn {
					return returnValue
				}

			} else {
				print.Message("Sorry, I can't pull the repository 😢", print.Error)
				print.Message("If there is a conflict, you can use the git cli to resolve it. \nOr you can rebase if you want\n	git pull -r "+remote.Name+" [branch] "+" \n	git push "+remote.Name+" [branch] ", print.None)
				os.Exit(1)
			}
		} else {
			print.Message("Sorry, I can't pull the repository 😢. I couldn't resolve this as a fast-forward merge and I can't use the git cli because it's not installed", print.Error)
			print.Message("Please install git (https://git-scm.com/downloads) and try again", print.None)
			os.Exit(1)
		}

	} else if err == git.NoErrAlreadyUpToDate || err == transport.ErrEmptyRemoteRepository || err == nil { // If there is nothing to pull, we push
		print.Message("Pull successful 🎉", print.Success)

		// If there is nothing to push, we exit
		// If there is another unknown error, we exit
		shouldReturn, returnValue := push(remote, path, profileLocal)
		if shouldReturn {
			return returnValue
		}
	} else { // If there is another unknown error, we exit
		exitOnError("Sorry, I can't pull the repository 😢", err)
	}
	return nil
}

func push(remote executor.Remote, path string, profileLocal profile.Profile) (bool, error) {
	spinnerPush := spinner.New(spinner.CharSets[9], 100)
	spinnerPush.Prefix = "Pushing the repository to " + remote.Name + " "
	spinnerPush.Start()
	err := executor.Push(path, remote.Name, profileLocal.Username, profileLocal.Password)
	spinnerPush.Stop()

	if err == git.NoErrAlreadyUpToDate || err == nil {
		print.Message("I've successfully synced your repository 🎉", print.Success)
		return true, nil
	} else {

		exitOnError("Sorry, I can't push the repository 😢", err)
	}
	return false, nil
}
