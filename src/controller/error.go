package controller

import (
	"fmt"
	"os"

	"errors"

	"github.com/julien040/gut/src/print"
)

var (
	ErrorReadInput = GutError{
		Message: "Sorry, I can't read your input",
		Code:    1,
		Err:     errors.New("error reading input"),
	}
	ErrorWorkingTreeNotClean = GutError{
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
	return fmt.Sprintf("https://gut.dev/error/%d", GutError.Code)
}

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
