package controller

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/prompt"
	"github.com/julien040/gut/src/telemetry"
)

func Telemetry(cmd *cobra.Command, args []string) {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}
	// Set arg to lowercase
	arg = strings.ToLower(arg)
	switch arg {
	case "on", "enable", "true", "1", "yes", "y", "o":
		telemetry.SetConsentState(true)
		print.Message("Telemetry is now enabled", print.Success)
	case "off", "disable", "false", "0", "no", "n", "f":
		telemetry.SetConsentState(false)
		print.Message("Telemetry is now disabled", print.Success)
	default:
		res, err := prompt.InputBool("Do you want to enable telemetry?", true)
		if err != nil {
			exitOnKnownError(errorReadInput, err)
		}
		telemetry.SetConsentState(res)
		if res {
			print.Message("Telemetry is now enabled", print.Success)
		} else {
			print.Message("Telemetry is now disabled", print.Success)
		}
	}

}
