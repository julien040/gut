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
	"github.com/julien040/gut/src/controller"

	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save [-e=[editor]] [-m=message] [-t=title] [files...]",
	Short: "Save (commit) your current work locally",
	Long: `Save (commit) your current work locally
To commit only some files, pass them as arguments to the command.
In case no files are passed as arguments, all files will be committed.

The -e flag allows you to specify the editor to use for writing the commit message.
If you want to use the default Git editor, you can also use the -e flag without any argument (e.g., git save -e).
To specify the editor, use the -e flag followed by the editor command (e.g. gut save -e="mate -w")`,
	Aliases: []string{"s", "commit"},
	Run:     controller.Save,
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringP("message", "m", "", "The commit message")
	saveCmd.Flags().StringP("title", "t", "", "The title of the commit")

	// https://github.com/spf13/pflag#setting-no-option-default-values-for-flags
	// To set the default value of a flag to an empty string, use the NoOptDefVal field.
	// -e return config
	// -e "mate -w" return mate -w
	// no flag return none
	saveCmd.Flags().StringP("editor", "e", "none", "The editor to use to write the commit message. Set -e to use the default git editor or specify one with -e=\"editor\"")
	saveCmd.Flag("editor").NoOptDefVal = "config"

}
