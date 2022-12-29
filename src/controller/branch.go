package controller

import (
	"github.com/julien040/gut/src/print"
	"github.com/spf13/cobra"
)

func Branch(cmd *cobra.Command, args []string) {
	print.Message("I'm in branch", "blue")
}

func BranchDelete(cmd *cobra.Command, args []string) {
	print.Message("I'm in branch delete", "blue")
}
