package controller

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

func validatorURL(val interface{}) error {
	url, ok := val.(string)
	if !ok {
		return errors.New("url is not a string")
	}
	if url == "" {
		return errors.New("url is empty")
	}
	isValid := checkURL(url)
	if !isValid {
		return errors.New("url is not valid")
	}
	return nil

}

func addRemote(path string, origin bool) (executor.Remote, error) {
	var qs []*survey.Question
	var answers struct {
		Name string
		Url  string
	}

	if !origin {
		qs = append(qs, &survey.Question{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name of the remote:",
			},
			Validate: survey.Required,
		})
	}
	qs = append(qs, &survey.Question{
		Name: "url",
		Prompt: &survey.Input{
			Message: "URL of the remote:",
		},
		Validate: validatorURL,
	})
	err := survey.Ask(qs, &answers)
	if err != nil {
		exitOnError("Sorry, I can't get your answer", err)
	}
	if origin {
		answers.Name = "origin"
	}
	err = executor.AddRemote(path, answers.Name, answers.Url)
	if err != nil {
		exitOnError("Sorry, I can't add the remote", err)
	}
	return executor.Remote{
		Name: answers.Name,
		Url:  answers.Url,
	}, nil

}

func chooseRemote(path string) (executor.Remote, error) {
	remote, err := executor.ListRemote(path)
	if err != nil {
		exitOnError("Sorry, I can't get the remote of the repository", err)
	}
	remoteName := make([]string, len(remote))
	for i, r := range remote {
		remoteName[i] = r.Name + " <" + r.Url + ">"
	}
	var remoteChoice int
	prompt := &survey.Select{
		Message: "Choose a remote:",
		Options: remoteName,
		Help:    "Choose a remote to use",
	}
	err = survey.AskOne(prompt, &remoteChoice)
	if err != nil {
		exitOnError("Sorry I can't get your answer", err)
	}
	return remote[remoteChoice], nil

}

func getRemote(path string) (executor.Remote, error) {
	remote, err := executor.ListRemote(path)
	if err != nil {
		exitOnError("Sorry, I can't get the remote of the repository", err)
	}
	lenRemote := len(remote)
	// Case no remote : We ask the user to add one
	if lenRemote == 0 {
		return addRemote(path, true)
	} else if lenRemote == 1 { // Case one remote : We return it
		return remote[0], nil
	} else { // Case multiple remote : We ask the user to choose one
		return chooseRemote(path)
	}

}

func Remote(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current directory", err)
	}
	checkIfGitRepoInitialized(wd)

	// List remote
	remote, err := executor.ListRemote(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the remote of the repository", err)
	}
	if len(remote) == 0 {
		print.Message("No remote found", print.Warning)
		res, err := prompt.InputBool("Do you want to add one ?", true)
		if err != nil {
			exitOnError("Sorry, I can't get your answer", err)
		}
		if res {
			remoteCreated, err := addRemote(wd, true)
			if err != nil {
				exitOnError("Sorry, I can't add the remote", err)
			}
			print.Message("I've successfully added the remote "+remoteCreated.Name+" <"+remoteCreated.Url+">", print.Success)

		}
	} else {
		if len(remote) == 1 {
			color.Black("Remote found:")
		} else {
			color.Black("Remotes found:")
		}
		for _, r := range remote {
			color.White("	%s <%s>", r.Name, color.BlackString(r.Url))
		}
	}

}

func RemoteAdd(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current directory", err)
	}
	checkIfGitRepoInitialized(wd)

	// List remote to check if we need to add origin
	remotes, err := executor.ListRemote(wd)
	if err != nil {
		exitOnError("Sorry, I can't get the remote of the repository", err)
	}
	// If no remote( len(remotes) == 0 ), it means that we need to add origin
	remote, err := addRemote(wd, len(remotes) == 0)
	if err != nil {
		exitOnError("Sorry, I can't add the remote", err)
	}

	print.Message("I've successfully added the remote "+remote.Name+" <"+remote.Url+">", print.Success)

}

func RemoteRemove(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get the current directory", err)
	}
	checkIfGitRepoInitialized(wd)

	var remoteDeleteName string
	if len(args) > 0 {
		// Check if the remote exists
		exists, err := executor.RemoteExists(wd, args[0])
		if err != nil {
			exitOnError("Sorry, I can't check if the remote exists", err)
		}
		if !exists {
			print.Message("The remote "+args[0]+" doesn't exist", print.Warning)
		} else {
			remoteDeleteName = args[0]
		}

	}

	if remoteDeleteName == "" {
		// List remotes
		remotes, err := executor.ListRemote(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the remote of the repository", err)
		}
		if len(remotes) == 0 {
			print.Message("No remote found", print.Warning)
			return
		}
		choosenRemote, err := chooseRemote(wd)
		if err != nil {
			exitOnError("Sorry, I can't get the remote to delete", err)
		}
		remoteDeleteName = choosenRemote.Name
	}

	res, err := prompt.InputBool("Are you sure you want to remove the remote "+remoteDeleteName+" ?", false)
	if err != nil {
		exitOnError("Sorry, I can't get your answer", err)
	}
	if !res {
		return
	}
	err = executor.RemoveRemote(wd, remoteDeleteName)
	if err != nil {
		exitOnError("Sorry, I can't remove the remote", err)
	}
	print.Message("I've successfully removed the remote "+remoteDeleteName, print.Success)

}
