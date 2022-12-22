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

func InputBool(message string) (bool, error) {
	res, err := InputWithValidation(message+" (Y/n)", "\nPlease enter either y or n",
		func(str string) bool {
			return str == "y" || str == "n" || str == "Y" || str == "N" || str == ""
		})
	if err != nil {
		return false, err
	}
	return res == "y" || res == "Y" || res == "", nil
}
