package controller

import (
	"github.com/julien040/gut/src/print"
	"github.com/spf13/cobra"
)

func Switch(cmd *cobra.Command, args []string) {
	print.Message("I'm in switch", "blue")
}
