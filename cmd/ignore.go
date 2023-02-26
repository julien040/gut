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

// ignoreCmd represents the ignore command
var ignoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "Download a .gitignore template",
	Long: `Download a .gitignore template from
https://github.com/toptal/gitignore
and add it to the current repository.
If .gitignore already exists, it will be updated.
`,
	Args: cobra.MaximumNArgs(1),

	Run:     controller.Ignore,
	Aliases: []string{"ig", "gitignore"},
}

var listTemplatesCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all available templates",
	Run:     controller.IgnoreList,
	Aliases: []string{"ls"},
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
	ignoreCmd.AddCommand(listTemplatesCmd)
}
