package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Get an input from the user using the specified editor.
// Similar to the git commit command that opens the editor to write the commit message
//
// Return the input as a string
func EditorGetText(editor string, defaultText string, wd string, fileTitle string) (string, error) {
	path := filepath.Join(wd, fileTitle)

	// Create or reset the file to the default text
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	file.WriteString(defaultText)
	file.Close()

	// Open the file with the editor
	// We have to split the editor string to get the command and the arguments
	commandName := strings.Split(editor, " ")[0]
	commandArgs := append(strings.Split(editor, " ")[1:], path)
	cmd := exec.Command(commandName, commandArgs...)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	// Read the file
	commitMessage, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Remove the file
	err = os.Remove(path)
	if err != nil {
		return "", err
	}

	return string(commitMessage), nil
}
