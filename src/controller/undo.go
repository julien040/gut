package controller

import (
	"github.com/julien040/gut/src/print"
	"github.com/spf13/cobra"
)

func Undo(cmd *cobra.Command, args []string) {
	print.Message("Uh oh, I can't do that yet 😢", print.Error)
}
