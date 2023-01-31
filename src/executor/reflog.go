package executor

import (
	"fmt"
	"strings"
)

type Reflog struct {
	NameToCheckout  string
	Action          string
	AdditionnalInfo string
	Date            string
}

// Parse the output of git reflog
func GetReflog(path string) (*[]Reflog, error) {
	output, err := gitReflog()
	if err != nil {
		return nil, err
	}

	var reflogs []Reflog

	var actionsAvailable = map[string]string{
		"pull":              "Pulled",
		"push":              "Pushed",
		"commit (initial):": "Initial commit",
		"commit (amend)":    "Modfication of a commit",
		"commit":            "Committed",
		"checkout":          "Switched",
		"reset":             "Reset",
		"revert":            "Reverted",
		"merge":             "Merged",
		"tag":               "Tagged",
	}

	// We have to use a slice because the order of the map is not guaranteed
	// This is ugly but it works
	// A refactor would be nice
	actionsSlice := []string{"pull", "push", "commit (initial):", "commit (amend)", "commit", "checkout", "reset", "revert", "merge", "tag"}

	// Split the output by line to get each reflog
	lineByLine := strings.Split(output, "\n")
	for i, line := range lineByLine {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Split the line by the separator specified in ./git.go:254
		splitedLine := strings.Split(line, "::::")
		// The first element is the date
		date := splitedLine[0]
		// All the other elements are the infos
		infos := strings.Join(splitedLine[1:], "")

		var additionnalInfo, action, nameToCheckout string

		// Check if the line contains an action we are looking for
		for _, actionName := range actionsSlice {

			if strings.HasPrefix(infos, actionName) {
				// If the action is found, we get the action name and the additionnal info
				action = actionsAvailable[actionName]
				additionnalInfo = strings.TrimPrefix(infos, actionName+":")

				// If the action is a pull or a push, it might contain credentials in the URL
				// We don't want to display them
				if actionName == "pull" || actionName == "push" {
					additionnalInfo = ""
				}

				break
			}
		}

		nameToCheckout = fmt.Sprintf("HEAD@{%d}", i)

		reflogs = append(reflogs, Reflog{
			NameToCheckout:  nameToCheckout,
			Action:          action,
			AdditionnalInfo: additionnalInfo,
			Date:            date,
		})

	}

	return &reflogs, nil

}
