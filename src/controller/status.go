package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
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

	checkIfGitRepoInitialized(wd)

	// Check if the head is detached
	detached, err := executor.IsDetachedHead(wd)
	if err != nil {
		exitOnError("Sorry, I can't check if the HEAD is detached 😢", err)
	}

	/* --------------------------- Get current branch --------------------------- */
	if detached {
		fmt.Println("HEAD is detached (lookout mode)")
	} else {
		branch, err := executor.GetCurrentBranch(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the current branch 😢", err)
		}
		fmt.Fprintf(color.Output, "On branch %s\n", color.HiGreenString(branch))
	}

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Prefix = "Getting status of files "
	spinner.Start()
	/* --------------------------- Get status of files -------------------------- */
	status, err := executor.GetStatus(wd)
	spinner.Stop()
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
