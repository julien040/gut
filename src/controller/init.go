package controller

import (
	"os"
	"path/filepath"

	"github.com/julien040/gut/src/executor"
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/spf13/cobra"
)

func Init(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		exitOnError("Oups, something went wrong while getting the current working directory", err)
	}
	// Profile selection is done before the git init to avoid an empty repository
	// If the profile selection fails, we don't want a git repo without any commits
	// Gut wouldn't be able to find the HEAD commit
	profile := selectProfile("", true)
	if executor.IsPathGitRepo(wd) {
		exitOnError("Oups, this directory is already a git repository. Delete the .git folder if you want to initialize a new repository", nil)
	}
	err = executor.Init(wd)
	if err != nil {
		exitOnError("Oups, something went wrong while initializing the repository", err)
	}

	// Ask for a .gitignore
	res, err := prompt.InputBool("Do you want to create a .gitignore file?", true)
	if err != nil {
		exitOnKnownError(errorReadInput, err)
	}
	if res {
		getGitignoreContentFromPath(wd)
		templates, err := fetchIgnoreList()
		if err != nil {
			exitOnError("Sorry, I couldn't fetch the list of available gitignore templates 😓", err)
		}
		template := selectGitignoreTemplate(templates)
		gitignoreTemplateContent := splitStringByNewLine(template.contents)
		gitignoreTemplateContent = append(gitignoreTemplateContent, ".gut")
		appendToGitignore(filepath.Join(wd, ".gitignore"), gitignoreTemplateContent, template.name)
	}

	associateProfileToPath(profile, wd)
	_, err = executor.Commit(wd, "🎉 Initial commit from Gut", nil)
	if err != nil {
		exitOnError("Oups, something went wrong while creating the first commit", err)
	}
	print.Message("Yeah, your repository is ready to go!", print.Success)

}
