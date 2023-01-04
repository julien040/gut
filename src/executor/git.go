/***********************************************************************************************************************
 *                 THIS FILE CONTAINS THE FUNCTIONS THAT EXECUTE GIT COMMANDS USING THE EXEC PACKAGE.                  *
 * BEFORE USING THESE FUNCTIONS, YOU SHOULD CHECK IF GIT IS INSTALLED ON THE SYSTEM USING THE ISGITINSTALLED FUNCTION. *
 *                     TRY FIRST TO USE THE GO-GIT PACKAGE, AND IF IT FAILS, USE THE EXEC PACKAGE.                     *
 ***********************************************************************************************************************/

package executor

import (
	"net/url"
	"os"
	"os/exec"
)

// Execute a command using the exec package and return the error if any
func runCommand(arg ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Dir = wd
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Check if git is installed on the system
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Execute git pull --rebase using the exec package
func GitPullRebase(username string, password string, remote string) error {
	// Parse the URL
	parsedURL, err := url.Parse(remote)
	if err != nil {
		return err
	}
	// Set the username and password
	parsedURL.User = url.UserPassword(username, password)

	// Execute the command
	err = runCommand("git", "pull", "--rebase", parsedURL.String())
	return err

}

// Execute git revert --no-edit using the exec package
func GitRevert(commitID string) error {
	return runCommand("git", "revert", "-m", "1", "--no-edit", commitID)
}