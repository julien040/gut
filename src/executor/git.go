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
// Print output only if it's stderr
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

// Execute a command using the exec package and return the error if any
// Print output only if it's stderr
// Return output as a string
func runCommandWithOutput(arg ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Dir = wd
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	return string(out), err
}

// Check if git is installed on the system
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Execute git pull using the exec package
func GitPull(username string, password string, remote string) error {
	// Parse the URL
	parsedURL, err := url.Parse(remote)
	if err != nil {
		return err
	}
	// Set the username and password
	parsedURL.User = url.UserPassword(username, password)

	// Execute the command
	err = runCommand("git", "pull", parsedURL.String())
	return err
}

// Set git config pull.rebase true only if not set to a value
func GitConfigPullRebaseTrueOnlyIfNotSet() error {

	// Check if the value is set
	out, err := runCommandWithOutput("git", "config", "pull.rebase")
	if err != nil {
		// If not set, Git will return exit status 1
		if err.Error() != "exit status 1" {
			return err
		}

	}
	// If the value is set, return
	if out != "" {
		return nil
	}

	// Set the value to true
	return runCommand("git", "config", "pull.rebase", "true")
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
	/*
		We use the commitID + ".." + "HEAD" syntax to revert the commit and all the commits after it.
		For example, if we have the following commits:
			- 1
			- 2
			- 3
			- 4

		And we want to revert the commit 2, we will use the following command:
			git revert -m 1 --no-edit 2..HEAD

		If we only use the commit id, it will revert change from this commit and not the commits after it.

	*/
	commitString := commitID + "..HEAD"
	return runCommand("git", "revert", "-m", "1", "--no-edit", commitString)
}

func GitRmCached() error {
	err := runCommand("git", "add", ".")
	if err != nil {
		return err
	}
	err = runCommand("git", "rm", "-r", "--ignore-unmatch", "--cached", ".")
	if err != nil {
		return err
	}
	return runCommand("git", "add", ".")
}

func GitResetAllHead() error {
	return runCommand("git", "reset", "--hard", "HEAD")
}

// I prefer to use checkout because git reset gave me some problems
// When used with a file in args, git reset was thinking it was a directory.
// Because of that, git reset was just exiting with an error
func GitCheckoutFileHead(file string) error {

	return runCommand("git", "checkout", "HEAD", "--", file)
}
