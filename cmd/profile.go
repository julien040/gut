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

	"github.com/julien040/gut/src/controller"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage your profiles",

	Run: controller.Profiles,
}

var profileAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new profile",
	Run:     controller.ProfilesAdd,
	Aliases: []string{"a", "new"},
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove [profile name]",
	Short: "Remove a profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile remove called")
	},
	Aliases: []string{"rm"},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all your profiles",
	Run:     controller.ProfilesList,
	Aliases: []string{"ls"},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profileListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
