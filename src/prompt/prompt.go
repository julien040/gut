package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/julien040/gut/src/print"
)

// InputLine prompts the user for an input and returns it
func InputLine(message string) (string, error) {
	fmt.Printf("%s ", message)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Remove the delimiter
	return strings.TrimSuffix(input, "\n"), nil
}

// Prompt the user for a string input and validate it
// If the input is not valid, the error message is displayed and the user is prompted again
// The validation function should return true if the input is valid
func InputWithValidation(message string, errorMessage string, validation func(string) bool) (string, error) {
	valid := false
	var res string
	for !valid {
		input, err := InputLine(message)
		if err == nil {
			valid = validation(input)
			if !valid {
				//Delete the last two lines
				print.Message(errorMessage, print.Error)
			} else {
				res = input
			}
		}
	}
	return res, nil
}

// InputBool prompts the user for a boolean input and returns it
//
// The default value is used if the user doesn't enter anything
// The message should not contain the "(y/n)" part
func InputBool(message string, defaultValue bool) (bool, error) {
	if defaultValue {
		message += " (Y/n): "
	} else {
		message += " (y/N): "
	}
	res, err := InputWithValidation(message, "\nPlease enter either y or n",
		func(str string) bool {
			return str == "y" || str == "n" || str == "Y" || str == "N" || str == ""
		})
	if err != nil {
		return false, err
	}
	if res == "" {
		return defaultValue, nil
	}
	return res == "y" || res == "Y", nil
}
