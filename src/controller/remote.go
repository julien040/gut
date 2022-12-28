package controller

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/julien040/gut/src/executor"
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
		exitOnError("We can't get your answers 😢", err)
	}
	if origin {
		answers.Name = "origin"
	}
	err = executor.AddRemote(path, answers.Name, answers.Url)
	if err != nil {
		exitOnError("We can't add the remote 😢", err)
	}
	return executor.Remote{
		Name: answers.Name,
		Url:  answers.Url,
	}, nil

}

func chooseRemote(path string) (executor.Remote, error) {
	remote, err := executor.ListRemote(path)
	if err != nil {
		exitOnError("We can't get the remote of the repository 😢", err)
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
		exitOnError("We can't get your choice 😢", err)
	}
	return remote[remoteChoice], nil

}

func getRemote(path string) (executor.Remote, error) {
	remote, err := executor.ListRemote(path)
	if err != nil {
		exitOnError("We can't get the remote of the repository 😢", err)
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
