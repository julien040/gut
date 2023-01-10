package controller

import (
	"fmt"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"

	"os"
	"path/filepath"
)

/*
I use the git CLI because I couldn't find any options to checkout only specific files using go-git
*/

func Undo(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't the current working directory", err)
	}
	checkIfGitRepoInitialized(wd)

	checkIfGitInstalled()

	if len(args) == 0 {
		res, err := prompt.InputBool("This will revert your working tree back to the state of the last commit. Any changes you made after the last commit will be lost. Are you sure you want to continue?", false)
		if err != nil {
			exitOnError("Sorry, I can't get your answer", err)
		}
		if !res {
			print.Message("Ok, I won't do anything", print.Info)
			return
		}
		err = executor.GitResetAllHead()
		if err != nil {
			exitOnError("Sorry, I can't revert your working tree to the last commit", err)
		}
		print.Message("I've successfully reverted your working tree to the last commit 🎉", print.Success)
		return
	} else {
		var resetFiles []string
		for _, val := range args {
			pathToCheck := filepath.Join(wd, val)
			// Check if path exits
			stat, err := os.Stat(pathToCheck)
			if err == nil {
				// Check if it's a file
				if !stat.IsDir() {
					resetFiles = append(resetFiles, val)
				}
			}

		}
		if len(resetFiles) == 0 {
			print.Message("Sorry, I can't find any file to undo", print.Error)
			return
		}
		res, err := prompt.InputBool(fmt.Sprintf("I will revert %d files to the last commit. Are you sure you want to continue?", len(resetFiles)), false)
		if err != nil {
			exitOnError("Sorry, I can't get your answer", err)
		}
		if !res {
			print.Message("Ok, I won't do anything", print.Info)
		}

		for _, val := range resetFiles {
			err = executor.GitCheckoutFileHead(val)
			if err == nil {
				print.Message(fmt.Sprintf("I've successfully reverted \"%s\" to the last commit 🎉", val), print.Success)
			} else {
				print.Message(fmt.Sprintf("Sorry, I can't revert \"%s\" to the last commit", val), print.Error)
			}

		}
		return

	}

}
