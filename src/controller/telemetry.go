package controller

import (
	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/print"
)

func Telemetry(cmd *cobra.Command, args []string) {
	print.Message("Telemetry is now disabled. Gut doesn't collect any data from you anymore.", print.Info)

}
