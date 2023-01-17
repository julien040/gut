package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func Merge(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory 😢", err)
	}
	checkIfGitRepoInitialized(wd)

	// Check if working directory is clean
	clean, err := executor.IsWorkTreeClean(wd)
	if err != nil {
		exitOnError("Sorry, I can't check if the working directory is clean 😢", err)
	}
	if !clean {
		print.Message("Sorry, but the working directory is not clean 😢", print.Error)
		res, err := prompt.InputBool("Would you like to save your changes?", true)
		if err != nil {
			exitOnError("Sorry, I can't get your answer 😢", err)
		}
		if res {
			// Run save command
			Save(cmd, args)
		} else {
			print.Message("See you later then 😊", print.Info)
			os.Exit(0)
		}

	}

	var branch string

	askForBranch := func() string {
		// Get the list of branches
		branches, err := executor.ListBranches(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the list of branches 😢", err)
		}
		// Remove the current branch from the list
		currentBranch, err := executor.GetCurrentBranch(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the current branch 😢", err)
		}
		var options []string
		for _, branch := range branches {
			if branch != currentBranch {
				options = append(options, branch)
			}
		}

		// Case only one branch so we can't merge
		if len(options) == 0 {
			exitOnError("Sorry, I can't find any branch to merge 😢", nil)
		}
		// Ask the user to choose a branch
		res, err := prompt.InputSelect("Choose a branch to merge into "+currentBranch, options)
		if err != nil {
			exitOnError("Sorry, I can't get your answer 😢", err)
		}

		return res
	}

	// Check if the user has specified a branch
	if len(args) > 0 {
		branch = args[0]
		// Check if the branch exists
		exists, err := executor.CheckIfBranchExists(wd, branch)
		if err != nil {
			exitOnError("Sorry, I can't check if the branch exists 😢", err)
		}
		// If not, ask him
		if !exists {
			print.Message(fmt.Sprintf("I'm sorry, but the branch %s doesn't exist 😢", branch), print.Error)
			branch = askForBranch()
		}
	} else {
		// If not, ask him
		branch = askForBranch()
	}

	currentBranch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the current branch 😢", err)
	}

	// Prompt the user to sync his changes
	// If the repo is not sync with the upstream, the pull request will not have the latest changes
	promptUserToSync := func() {
		res, err := prompt.InputBool("Would you like to sync your changes?", true)
		if err != nil {
			exitOnError("Sorry, I can't get your answer 😢", err)
		}
		if res {
			// Run sync command
			Sync(cmd, args)
		}
	}

	// Check if the repo is using only one origin. If yes, we act accordingly to each platform
	// If not, returns an empty string. In this case, we use the default merge method
	platform, remoteURL, err := getPlatformUsed(wd)
	if err != nil {
		print.Message("I can't know if you use Github, Gitlab or Bitbucket 😢. I will use the default merge method. Check the message below to know why I can't know.", print.Warning)
		print.Message(err.Error(), print.Warning)
	}

	// Return the URL to open a pull request on the platform
	switch platform {
	case "github.com":
		githubURL := remoteURL + "/compare/" + currentBranch + "..." + branch + "?quick_pull=1"
		promptUserToSync()
		color.Black("To merge %s into %s, I recommend opening a pull request on GitHub. Open the following URL in your browser:\n%s\n", branch, currentBranch, color.WhiteString(githubURL))
		openInBrowser(githubURL)
	case "gitlab.com":
		gitlabURL := remoteURL + "/merge_requests/new?merge_request%5Bsource_branch%5D=" + branch + "&merge_request%5Btarget_branch%5D=" + currentBranch
		promptUserToSync()
		color.Black("To merge %s into %s, I recommend opening a merge request on GitLab. Open the following URL in your browser:\n%s\n", branch, currentBranch, color.WhiteString(gitlabURL))
		openInBrowser(gitlabURL)
	case "bitbucket.org":
		bitbucketURL := remoteURL + "/pull-requests/new?source=" + branch + "&dest=" + currentBranch
		promptUserToSync()
		color.Black("To merge %s into %s, I recommend opening a pull request on Bitbucket. Open the following URL in your browser:\n%s\n", branch, currentBranch, color.WhiteString(bitbucketURL))
		openInBrowser(bitbucketURL)
	default:
		isGitInstalled := executor.IsGitInstalled()
		if isGitInstalled {
			color.Black("To merge %s into %s, I recommend using the following command:\n%s\n", branch, currentBranch, color.WhiteString("git merge "+branch))
		} else {
			exitOnError("You are using a platform that I don't support yet and you don't have git installed. I can't help you sorry 😢. Please install git (https://git-scm.com/downloads) and try again.", nil)
		}

	}

}

// Get the remote host of a git repository
//
// Response:
// The host of the remote | the remote URL without the .git suffix | an error if there is one
//
// ⚠️ Note: If there is more than one remote, or none, we can't know so we return an empty string
func getPlatformUsed(path string) (string, string, error) {

	removeDotGit := func(url string) string {
		if strings.HasSuffix(url, ".git") {
			return url[:len(url)-4]
		}
		return url
	}

	// Get the remotes
	remotes, err := executor.ListRemote(path)
	if err != nil {
		return "", "", err
	}
	// We can only know the platform if there is only one remote
	if len(remotes) == 1 {
		return getHost(remotes[0].Url), removeDotGit(remotes[0].Url), nil
	}
	// If there is more, or none, we can't know so we return an empty string
	return "", "", nil
}
