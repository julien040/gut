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

	"github.com/julien040/gut/src/controller"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gut",
	Short: "Simplify your Git journey with Gut",
	Long: `
	Gut is a powerful command-line interface (CLI) designed to make Git easier to use.
	Effortlessly version control your projects with Gut.`,
	Run: controller.Root,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Run the command called by the user
	err := rootCmd.Execute()

	// Exit with error code 1 if an error occured
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print the version of Gut and the release date")

}
