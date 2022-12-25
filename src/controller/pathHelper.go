/* -------------------------------------------------------------------------- */
/*      This file contains function that can be used by multiple commands     */
/* -------------------------------------------------------------------------- */

package controller

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/julien040/gut/src/prompt"
)

// Return an absolute path from the user input
func getAbsPathFromInput(repoName string, str string) string {
	res := ""
	wd, err := os.Getwd()
	if str == "." {
		res = wd
		if err != nil {
			exitOnError("For some reason, we couldn't get your input 😓", err)
		}
	} else if str == "" {
		res = filepath.Join(wd, repoName)
	} else if filepath.IsAbs(str) {
		res = str
	} else {
		// If the path is not absolute, we assume it's relative to the current directory
		res, err = filepath.Abs(str)
		if err != nil {
			exitOnError("We can't parse the path you entered. Make sure it's a valid path", err)
		}
	}
	return res
}

// Ask the user for a path and return the absolute path
func askForPath(repoName string, message string) string {
	str, err := prompt.InputLine(message +
		color.HiBlackString("\n . for the current directory | leave empty for a folder named "+repoName+" | or enter a path (e.g. /home/user/repo) \n"))
	if err != nil {
		exitOnError("For some reason, we couldn't get your input 😓", err)
	}
	return getAbsPathFromInput(repoName, str)
}

// Check if the path exists and return true if it does
//
// Path should be absolute
func checkIfPathExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Ask the user if he wants to create the path he entered
func promptUserForMakingPath(path string) bool {
	shouldCreate, err := prompt.InputBool("The path you entered doesn't exist. Do you want to create it? \n "+color.HiBlackString(path), true)
	if err != nil {
		exitOnError("For some reason, we couldn't get your input 😓", err)
	}
	if shouldCreate {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			exitOnError("We couldn't create the path you entered 😓. Please make sure you have the right permissions", err)
		}
	}
	return shouldCreate
}

// Returns true if the directory is empty
func isDirectoryEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == nil {
		return false, nil
	}
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
