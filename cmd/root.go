/*
Copyright © 2022 Julien CAGNIART

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/julien040/gut/src/prompt"
	"github.com/julien040/gut/src/telemetry"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gut",
	Short: "Show the status of your git repo",
	Long: `
Gut is a simple git client.
It allows you to manage your git repositories in a simple way.
Our goal is not to replace git, but to make it easier to use for the most common tasks.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// Prompt for telemetry consent
	promptForTelemetry()

	// Run the command called by the user
	err := rootCmd.Execute()

	// Log the command called by the user
	if len(os.Args) > 1 {
		telemetry.LogCommand(os.Args[1])
	} else {
		telemetry.LogCommand("gut")
	}

	// Exit with error code 1 if an error occured
	if err != nil {
		os.Exit(1)
	}

}

func promptForTelemetry() {
	consentStateKnown := telemetry.IsConsentStateKnown()
	if !consentStateKnown {
		res, err := prompt.InputBool("Do you want to help us improve gut by sending anonymous usage data?", true)
		if err != nil {
			return
		}
		telemetry.SetConsentState(res)
	}

}

func init() {

}
