package controller

import (
	"fmt"
	"os"

	"errors"

	"github.com/fatih/color"
)

var (
	errorReadInput = GutError{
		Message: "This is embarrassing, I can't read your input",
		Code:    1,
		Err:     errors.New("error reading input"),
	}
	errorWorkingTreeNotClean = GutError{
		Message: "I can't continue further because the working tree is not clean",
		Code:    2,
		Err:     errors.New("working tree not clean"),
	}
)

type GutError struct {
	Message string
	Code    int
	Err     error
}

// An helper function to exit the program if an error occurs
func exitOnError(str string, err error) {
	// Print the error message to stderr

	// Print a new line
	fmt.Fprintln(os.Stderr, "")
	if str != "" {
		// Print the error message in red
		fmt.Fprintf(os.Stderr, "%s \n", color.RedString(str))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s \n", err.Error())
		fmt.Fprintf(os.Stderr, "%s\n", "If this error persists, please open an issue on GitHub: https://github.com/julien040/gut/issues/new")
	}

	os.Exit(1)
}

func getLinkForError(GutError GutError) string {
	return fmt.Sprintf("https://gut-cli.dev/error/%d", GutError.Code)
}

// Same as exitOnError but with a GutError
//
// # This function will print the error message and the link to the error page
//
// This should be used when the error is known and can be resolved by the user
func exitOnKnownError(typeOfError GutError, err error) {
	// Print the error message to stderr
	fmt.Fprintln(os.Stderr, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, color.RedString("%s (error: %s)\n"), color.RedString(typeOfError.Message), color.RedString(err.Error()))
		/* print.Message("%s (error: %s)", print.Error, typeOfError.Message, err.Error()) */
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", color.RedString(typeOfError.Message))
		/* print.Message("%s", print.Error, typeOfError.Message) */
	}
	fmt.Fprintf(os.Stderr, "To resolve this issue, please follow the instructions on this page: %s\n", getLinkForError(typeOfError))
	/* print.Message("To resolve this issue, please follow the instructions on this page: %s", print.Optional, getLinkForError(typeOfError)) */

	os.Exit(1)
}
