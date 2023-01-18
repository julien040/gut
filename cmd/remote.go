/*
Copyright © 2023 Julien CAGNIART

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

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:     "remote",
	Short:   "List remote repositories",
	Aliases: []string{"r", "remotes", "remote ls"},
	Run:     controller.Remote,
}

var remoteAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a remote repository",
	Aliases: []string{"a", "new", "create"},
	Run:     controller.RemoteAdd,
}

var remoteRemoveCmd = &cobra.Command{
	Use:     "rm [remote name]",
	Short:   "Remove a remote repository",
	Aliases: []string{"remove", "delete", "del", "d"},
	Run:     controller.RemoteRemove,
	Args:    cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)
}
