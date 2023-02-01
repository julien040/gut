package controller

import (
	"os"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func Branch(cmd *cobra.Command, args []string) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't find the current working directory", err)
	}
	// Check if the current directory is a git repository and exit if not
	checkIfGitRepoInitialized(wd)

	// List all branches
	branches, err := executor.ListBranches(wd)
	if err != nil {
		exitOnError("I can't list the branches", err)
	}

	// Get the current branch
	currentBranch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("I can't get the current branch", err)
	}

	// Print the branches
	print.Message("Here is the list of all branches:", "none")
	for _, branch := range branches {
		if branch == currentBranch {
			print.Message("* "+branch, "green")
		} else {
			print.Message("  "+branch, "none")
		}
	}

}

func BranchDelete(cmd *cobra.Command, args []string) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't find the current working directory", err)
	}
	// Check if the current directory is a git repository and exit if not
	checkIfGitRepoInitialized(wd)

	// Check if there are uncommitted changes
	clean, err := executor.IsWorkTreeClean(wd)
	if err != nil {
		exitOnError("I can't check if there are uncommitted changes", err)
	}
	if !clean {
		res, err := prompt.InputBool("Uh oh, there are uncommitted changes. They might be lost if you delete the branch. Do you want to continue?", false)
		if err != nil {
			exitOnError("That's a shame, I can't get your answer", err)
		}
		if !res {
			print.Message("Okay, I won't delete the branch", print.Info)
			return
		}
	}

	// Input the branch name
	var branchName string
	if len(args) == 0 {
		res, err := prompt.InputLine("Branch to delete: ")
		if err != nil {
			exitOnError("Sorry, I can't get your answer", err)
		}
		branchName = res
	} else {
		branchName = args[0]
	}

	// Check if the branch exists
	exists, err := executor.CheckIfBranchExists(wd, branchName)
	if err != nil {
		exitOnError("I can't check if the branch exists", err)
	}

	if !exists {
		print.Message("The branch doesn't exist. I can't delete it", print.Error)
		return
	}

	// Check if the branch is the current branch
	currentBranch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("I can't get the current branch", err)
	}

	if branchName == currentBranch {
		print.Message("You can't delete the current branch. I recommend you to switch to another branch first", print.Error)
		return
	}

	res, err := prompt.InputBool("Are you sure you want to delete the local branch "+branchName+"?", false)
	if err != nil {
		exitOnError("Sorry, I can't get your answer", err)
	}
	if !res {
		print.Message("Okay, I won't delete the branch", print.Info)
		return
	}

	// Delete the branch
	err = executor.DeleteBranch(wd, branchName)
	if err != nil {
		exitOnError("I can't delete the branch", err)
	}

	// Delete the remote branch
	installed := executor.IsGitInstalled()
	if installed {
		// Check if there is one remote
		remotes, err := executor.ListRemote(wd)
		if err != nil {
			exitOnError("I can't list the remotes", err)
		}
		if len(remotes) == 1 {
			remote := remotes[0]
			res, err := prompt.InputBool("Do you want to delete the remote branch "+remote.Name+"/"+branchName+"?", false)
			if err != nil {
				exitOnError("Sorry, I can't get your answer", err)
			}
			if res {
				err := executor.GitDeleteRemoteBranch(remote.Name, branchName)
				if err != nil {
					exitOnError("I can't delete the remote branch", err)
				}
			}

		}
	}

	print.Message("The branch has been deleted", print.Success)
}
