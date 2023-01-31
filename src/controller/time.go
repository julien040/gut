package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

func Time(cmd *cobra.Command, args []string) {

	wd := getWorkingDir()

	checkIfGitRepoInitialized(wd)

	checkIfGitInstalled()

	reflogs, err := executor.GetReflog(wd)
	if err != nil {
		exitOnError("Sorry, I couldn't get the reflog", err)
	}

	var humanReadableReflogs []string
	for _, reflog := range *reflogs {
		dateParsed, err := time.Parse("2006-01-02 15:04:05 -0700", reflog.Date)
		if err != nil {
			print.Message("I couldn't parse the date of the reflog: %s", print.Warning, err)
			continue
		}

		// Trim additional info
		reflog.AdditionnalInfo = strings.TrimSpace(reflog.AdditionnalInfo)

		options := fmt.Sprintf("%s (%s) on %s", color.BlueString(reflog.Action), color.BlackString(reflog.AdditionnalInfo), color.WhiteString(dateParsed.Format("Mon Jan 2 15:04:05")))

		// Case there is no additionnal info
		if reflog.AdditionnalInfo == "" {
			options = fmt.Sprintf("%s on %s", color.BlueString(reflog.Action), color.WhiteString(dateParsed.Format("Mon Jan 2 15:04:05")))
		}

		humanReadableReflogs = append(humanReadableReflogs, options)
	}

	// Ask the user which reflog he wants to checkout
	var reflogToCheckout int
	selectPrompt := &survey.Select{
		Message:  "Which change do you want to look at?",
		Options:  humanReadableReflogs,
		PageSize: 10,
		Help:     "This is the history of changes in your repository. You can select one to see the state of your repository at this time.",
	}
	err = survey.AskOne(selectPrompt, &reflogToCheckout)
	if err != nil {
		exitOnError("Sorry, I can't get your answer 😢", err)
	}
	res, err := prompt.InputBool("Are you sure you want to lookout at this time?", true)
	if err != nil {
		exitOnError("Sorry, I can't get your answer 😢", err)
	}
	if !res {
		return
	}
	err = executor.GitCheckout((*reflogs)[reflogToCheckout].NameToCheckout)
	if err != nil {
		exitOnError("Sorry, I can't checkout this reflog 😢", err)
	}
	print.Message("I've successfully switched to this time!", print.Success)
	print.Message("To go back to a branch, use:\n	gut switch <branch", print.Optional)
	print.Message("To create a branch from this commit, use:\n	gut switch [new branch name]", print.Optional)

}
