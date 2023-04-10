package controller

import (
	"github.com/julien040/gut/src/print"
	"github.com/julien040/gut/src/telemetry"
	"github.com/spf13/cobra"
)

func Root(cmd *cobra.Command, args []string) {
	// Check if the version flag is set
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		exitOnError("Sorry, I can't retrieve the version flag 😓", err)
	}

	// Retrieve Gut version and release date
	gutVersion := telemetry.GetBuildInfo()

	// Case gut is called with -v
	if version {
		// When gut is built from source, the version is set to "dev"
		// We handle this case to avoid printing "gut vdev (1970-01-01)"
		if gutVersion == "dev" {
			print.Message("You are running a built from source version of Gut", print.None)
		} else {
			// Print Gut version and release date
			print.Message("gut v%s", print.None, gutVersion)
		}
	} else { // Case gut is called without -v
		// We show the help message
		cmd.Help()
	}

}
