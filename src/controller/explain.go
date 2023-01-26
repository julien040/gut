package controller

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
)

func Explain(cmd *cobra.Command, args []string) {
	var arg string
	if len(args) == 0 {
		// Create a list of all the commands
		commands := []string{}
		for command := range explainMessage {
			commands = append(commands, command)
		}

		// Sort the list
		sort.Strings(commands)

		// Ask the user to choose a command
		res, err := prompt.InputSelect("Choose a command", commands)
		if err != nil {
			exitOnError("Sorry, I can't get your answer 😢", err)
		}
		arg = res

	} else {
		arg = args[0]
	}
	print.Message(findExplainMessageFromArg(arg), print.None)

}

func findExplainMessageFromArg(arg string) string {
	// Try to find the message in the map
	if message, ok := explainMessage[arg]; ok {
		return message
	}
	return "Sorry, I don't this command 😢"
}

var explainMessage = map[string]string{
	"branch": "Sorry I don't know how to explain this command yet",

	"clone": `To run this command with git, you would do:
	$ git clone <url> <directory>
	`,

	"diff": `To run this command with git, you would do:
	$ git diff --staged`,

	"explain": "42",

	"fix": "Sorry I don't know how to explain this command yet",

	"history": `To run this command with git, you would do:
	$ git log --oneline`,

	"ignore": "There is no git command for this",

	"init": `To run this command with git, you would do:
	$ git init`,

	"merge": `To run this command with git, you would do:
	$ git merge <branch>`,

	"profile": "There is no git command for this",

	"remote": "Sorry I don't know how to explain this command yet",

	"revert": `To run this command with git, you would do:
	$ git revert <commit>..HEAD`,

	"save": `To run this command with git, you would do:
	$ git commit -m "<message>"`,

	"status": `To run this command with git, you would do:
	$ git status`,
	"switch": `To run this command with git, you would do:
	$ git checkout <branch>`,
	"sync": "Sorry I don't know how to explain this command yet",
	"undo": `To run this command with git, you would do:
	$ git reset HEAD --hard`,
	"whereami": `To run this command with git, you would do:
	$ git rev-parse HEAD && git rev-parse --abbrev-ref HEAD`,
}
