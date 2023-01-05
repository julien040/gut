package controller

import (
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
		exitOnError("We can't get the current working directory 😢", err)
	}
	remote, err := getRemote(wd)
	if err != nil {
		exitOnError("We can't get the remote of the repository 😢", err)
	}
	syncRepo(wd, remote, false)
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
		print.Message("Uh oh, there is a conflict 😢. Currently, Gut doesn't support conflict resolution. Please resolve the conflict manually.", print.Error)
		print.Message("You can use the git cli to resolve the conflict. \n	git pull -r "+remote.Name+" \n	git push "+remote.Name, print.None)
		os.Exit(1)

	} else if err == git.NoErrAlreadyUpToDate || err == nil { // If there is nothing to pull or if there is no error, we push the repository
		print.Message("Pull successful 🎉", print.Success)

		spinnerPush := spinner.New(spinner.CharSets[9], 100)
		spinnerPush.Prefix = "Pushing the repository to " + remote.Name + " "
		spinnerPush.Start()
		err := executor.Push(path, remote.Name, profileLocal.Username, profileLocal.Password)
		spinnerPush.Stop()

		if err == git.NoErrAlreadyUpToDate || err == nil { // If there is nothing to push, we exit
			print.Message("The repository has been synced with "+remote.Name+" 🎉", print.Success)
			return nil
		} else {
			// If there is another unknown error, we exit
			exitOnError("We can't push the repository 😢", err)
		}
	} else { // If there is another unknown error, we exit
		exitOnError("We can't pull the repository 😢", err)
	}
	return nil
}
