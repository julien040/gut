package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

type ignoreTemplate struct {
	name     string
	contents string
}

const gitignoreApiURL = "https://www.toptal.com/developers/gitignore/api/list?format=json"

// Fetch the list of available gitignore templates from constant gitignoreApiURL
func fetchIgnoreList() ([]ignoreTemplate, error) {
	// Fetch the list of available gitignore templates
	print.Message("Fetching the list of available gitignore templates...", "info")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	resp, err := http.Get(gitignoreApiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var bodyAPI map[string]interface{}

	var templates []ignoreTemplate
	err = json.Unmarshal(body, &bodyAPI)
	if err != nil {
		return nil, err
	}
	for _, v := range bodyAPI {
		// Check if it's a map
		if template, ok := v.(map[string]interface{}); ok {
			templates = append(templates, ignoreTemplate{
				name:     template["name"].(string),
				contents: template["contents"].(string),
			})
		}
	}
	s.Stop()
	return templates, nil

}

// Ask the user to select a gitignore template
func selectGitignoreTemplate(templates []ignoreTemplate) ignoreTemplate {
	// Ask the user to select a gitignore template
	var templateNames []string
	for _, template := range templates {
		templateNames = append(templateNames, template.name)
	}
	var selectedTemplate int
	prompt := &survey.Select{
		Message: "Select a gitignore template:",
		Options: templateNames,
	}
	err := survey.AskOne(prompt, &selectedTemplate)
	if err != nil {
		exitOnKnownError(errorReadInput, err)
	}
	return templates[selectedTemplate]
}

func getGitignoreContentFromPath(path string) string {
	// Check if .gitignore already exists
	gitignorePath := filepath.Join(path, ".gitignore")
	file, err := os.Open(gitignorePath)
	if os.IsNotExist(err) {
		// Create the file
		file, err := os.Create(gitignorePath)
		if err != nil {
			exitOnError("Sorry, I can't create the .gitignore file 😓. Please, can you check if I have the right permissions?", err)
		}
		defer file.Close()
		return ""
	} else if err != nil {
		defer file.Close()
		exitOnError("Sorry, I can't open the .gitignore file 😓. Please, can you check if I have the right permissions?", err)

	} else {
		// Read the file
		content, err := io.ReadAll(file)
		if err != nil {
			exitOnError("Sorry, I can't read the .gitignore file 😓. Please, can you check if I have the right permissions?", err)
		}
		defer file.Close()
		return string(content)

	}

	return ""

}

func splitStringByNewLine(str string) []string {
	return strings.Split(str, "\n")
}

// Returns a list of elements that are in b and not in a
func Difference(a, b []string) []string {
	// https://stackoverflow.com/a/45428032/15573415

	// Create a map of all elements in a.
	m := make(map[string]bool)
	for _, x := range a {
		m[x] = true
	}
	var diff []string

	// Loop over b and check if each element is in the map.
	// If not, it means it's not in a and we append it to the diff.
	for _, x := range b {
		if !m[x] {
			// Check if the line is not empty and not a comment
			if x != "" && x[:1] != "#" {
				diff = append(diff, x)
			}
		}
	}
	return diff

}

// Append a list of excluded files to the .gitignore file
func appendToGitignore(path string, content []string, templateName string) {
	// Open the file in append mode
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		exitOnError("Sorry, I can't open the .gitignore file 😓. Please, can you check if I have the right permissions?", err)
	}
	defer file.Close()

	// Add a comment to the file
	_, err = file.WriteString("\n# " + templateName + " template downloaded with gut\n")
	if err != nil {
		exitOnError("Sorry, I can't write to the .gitignore file 😓. Please, can you check if I have the right permissions?", err)
	}

	// Append the content to the file
	for _, line := range content {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			exitOnError("Sorry, I can't write to the .gitignore file 😓. Please, can you check if I have the right permissions?", err)
		}
	}
}

// Cmd to download a gitignore template
func Ignore(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Sorry, I can't get your current working directory 😓", err)
	}
	// Get the local .gitignore content
	gitignoreContent := splitStringByNewLine(getGitignoreContentFromPath(wd))

	templates, err := fetchIgnoreList()
	if err != nil {
		exitOnError("Sorry, I couldn't fetch the list of available gitignore templates 😓", err)
	}
	template := selectGitignoreTemplate(templates)
	gitignoreTemplateContent := splitStringByNewLine(template.contents)
	// Get the difference between the local .gitignore and the gitignore template
	diff := Difference(gitignoreContent, gitignoreTemplateContent)
	// Append the difference to the local .gitignore
	if len(diff) == 0 {
		print.Message("Your .gitignore file is already up to date with the "+template.name+" template 😎", print.Success)
		return
	} else {
		appendToGitignore(filepath.Join(wd, ".gitignore"), diff, template.name)

		print.Message("I've updated your .gitignore file with the "+template.name+" template 🎉", print.Success)

		/*
			Git will continue to track the files that are already tracked by git even if they are added to the .gitignore file.
			To untrack the files, we need to run the command `git rm -r --cached .`
			We prompt the user if they want us to run the command for them.
		*/
		installed := executor.IsGitInstalled() // We prompt only if git is installed

		// Check if a repo is initialized
		repoInitialized := executor.IsPathGitRepo(wd)

		if installed && repoInitialized {
			print.Message("If you plan to use the git CLI, you might want to untrack the files that are already tracked by git.", print.Info)
			res, err := prompt.InputBool("Do you want me to run the command for you?", true)
			if err != nil {
				exitOnKnownError(errorReadInput, err)
			}
			if res {
				err := executor.GitRmCached()
				if err != nil {
					exitOnError("Sorry, I couldn't untrack the files 😓", err)
				} else {
					print.Message("I've untracked the files for you 🎉", print.Success)
				}

			}
		}
	}

}

// Cmd to list all available gitignore templates
func IgnoreList(cmd *cobra.Command, args []string) {
	templates, err := fetchIgnoreList()
	if err != nil {
		exitOnError("Sorry, I couldn't fetch the list of available gitignore templates 😓", err)
	}
	print.Message("Available gitignore templates:", print.Info)
	temp := selectGitignoreTemplate(templates)
	print.Message("Here's the content of the "+temp.name+" template:", print.Info)
	fmt.Println(temp.contents)
}
