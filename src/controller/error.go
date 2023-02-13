package controller

import (
	"fmt"
	"os"

	"errors"

	"github.com/julien040/gut/src/print"
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
	fmt.Println("")
	if str != "" {
		print.Message(str, print.Error)
	}
	if err != nil {
		fmt.Printf("Error message: %s\n", err)
		print.Message("If this error persists, please open an issue on GitHub: https://github.com/julien040/gut/issues/new", print.None)
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
	fmt.Println("")
	if err != nil {
		print.Message("%s (error: %s)", print.Error, typeOfError.Message, err.Error())
	} else {
		print.Message("%s", print.Error, typeOfError.Message)
	}
	print.Message("To resolve this issue, please follow the instructions on this page: %s", print.Optional, getLinkForError(typeOfError))

	os.Exit(1)
}
