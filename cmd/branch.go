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

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:     "branch",
	Short:   "List all branches",
	Long:    `List all branches of the repository`,
	Run:     controller.Branch,
	Aliases: []string{"b", "br", "branch ls", "branch list", "b ls", "br ls"},
}

var branchAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Alias of gut switch",
	Run:     controller.Switch,
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"new"},
}

var branchDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a branch",
	Long:    `Delete a branch from the repository`,
	Run:     controller.BranchDelete,
	Aliases: []string{"del", "remove", "rm"},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchAddCmd)
	branchCmd.AddCommand(branchDeleteCmd)

}
