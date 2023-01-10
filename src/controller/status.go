package controller

import (
	"fmt"
	"os"

	"github.com/julien040/gut/src/executor"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Status(cmd *cobra.Command, args []string) {
	/* -------------------------- Get current directory ------------------------- */
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current working directory 😢", err)
	}

	/* --------------------------- Get current branch --------------------------- */
	branch, err := executor.GetCurrentBranch(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the current branch 😢", err)
	}
	fmt.Printf("On branch %s\n", color.HiGreenString(branch))

	/* --------------------------- Get status of files -------------------------- */
	status, err := executor.GetStatus(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the status of the repository 😢", err)
	}
	if len(status) == 0 {
		fmt.Println("No changes to commit")
		return
	} else {
		fmt.Println("\nChanges to be committed:")

		for _, file := range status {
			if file.Status == "M" {
				color.Yellow("\tModified: %s", color.HiBlackString(file.Path))
			} else if file.Status == "A" {
				color.Green("\tNew file: %s", color.HiBlackString(file.Path))
			} else if file.Status == "D" {
				color.Red("\tDeleted: %s", color.HiBlackString(file.Path))
			} else if file.Status == "R" {
				color.White("\tRenamed: %s", color.HiBlackString(file.Path))
			} else if file.Status == "C" {
				color.White("\tCopied: %s", color.HiBlackString(file.Path))
			} else if file.Status == "?" {
				color.Cyan("\tUntracked: %s", color.HiBlackString(file.Path))
			} else if file.Status == "U" {
				color.Red("\tUpdated but unmerged: %s", color.HiBlackString(file.Path))
			}
		}
	}

}
