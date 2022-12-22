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
	"fmt"

	"github.com/spf13/cobra"
)

// ignoreCmd represents the ignore command
var ignoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "Download a .gitignore template",
	Long: `Download a .gitignore template from
https://github.com/github/gitignore
and add it to the current repository.
If .gitignore already exists, it will be updated or overwritten.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ignore called")
	},
}

var ignoreAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a file or a directory to the .gitignore file",
	Long:  `Add a file or a directory to the .gitignore file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ignore add called")
	},
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long:  `List all available templates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ignore list called")
	},
	Aliases: []string{"ls"},
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
	ignoreCmd.AddCommand(ignoreAddCmd)
	ignoreCmd.AddCommand(listTemplatesCmd)
	ignoreAddCmd.Flags().BoolP("force", "f", false, "Overwrite the .gitignore file if it already exists")
	ignoreAddCmd.Flags().BoolP("update", "u", false, "Update the .gitignore file if it already exists")
	ignoreAddCmd.Flags().String("template", "t", "Download a .gitignore template from github.com/github/gitignore")
	ignoreAddCmd.Flags().String("filename", "f", "Name of the file to add to the .gitignore file")
	ignoreAddCmd.MarkFlagFilename("filename")
	ignoreAddCmd.Flags().String("directory", "d", "Name of the directory to add to the .gitignore file")
	ignoreAddCmd.MarkFlagDirname("directory")
}
