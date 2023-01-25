package controller

import (
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
		//
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
	"branch":   "This feature is not implemented yet",
	"clone":    "This feature is not implemented yet",
	"diff":     "This feature is not implemented yet",
	"explain":  "This feature is not implemented yet",
	"fix":      "This feature is not implemented yet",
	"history":  "This feature is not implemented yet",
	"ignore":   "This feature is not implemented yet",
	"init":     "This feature is not implemented yet",
	"merge":    "This feature is not implemented yet",
	"profile":  "This feature is not implemented yet",
	"remote":   "This feature is not implemented yet",
	"revert":   "This feature is not implemented yet",
	"save":     "This feature is not implemented yet",
	"status":   "This feature is not implemented yet",
	"switch":   "This feature is not implemented yet",
	"sync":     "This feature is not implemented yet",
	"undo":     "This feature is not implemented yet",
	"whereami": "This feature is not implemented yet",
}
