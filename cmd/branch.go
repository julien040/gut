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

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List all branches",
	Long:  `List all branches of the repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("branch called")
	},
	Aliases: []string{"b", "br", "branch ls", "branch list", "b ls", "br ls"},
}

var branchAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new branch",
	Long:  `Add a new branch to the repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("branch add called")
	},
	Aliases: []string{"new"},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a branch",
	Long:  `Delete a branch from the repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("branch delete called")
	},
	Aliases: []string{"del", "remove", "rm"},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchAddCmd)
	branchCmd.AddCommand(branchDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
