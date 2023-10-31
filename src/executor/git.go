/***********************************************************************************************************************
 *                 THIS FILE CONTAINS THE FUNCTIONS THAT EXECUTE GIT COMMANDS USING THE EXEC PACKAGE.                  *
 * BEFORE USING THESE FUNCTIONS, YOU SHOULD CHECK IF GIT IS INSTALLED ON THE SYSTEM USING THE ISGITINSTALLED FUNCTION. *
 *                     TRY FIRST TO USE THE GO-GIT PACKAGE, AND IF IT FAILS, USE THE EXEC PACKAGE.                     *
 ***********************************************************************************************************************/

package executor

import (
	"bufio"
	"io"
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

// Iframe a command using the exec package and return the error if any
func runCommandWithStdin(arg ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)

	cmd := exec.Command(arg[0], arg[1:]...)
	cmd.Dir = wd
	cmd.Stderr = os.Stderr

	read, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	// Pipe the reader to the command
	go func() {
		defer read.Close()
		io.Copy(read, reader)
	}()

	// Same for stdout
	cmd.Stdout = os.Stdout

	return cmd.Run()
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

// Check if the remote contains the hash
// Used to check if the commit is already pushed to the remote
func GitRemoteContainsHash(hash string) bool {
	output, err := runCommandWithOutput("git", "branch", "-r", "--contains", hash)
	if err != nil {
		return true
	}
	return output != ""
}

// Fetch all the branches from the remote
func GitFetchAll() error {
	return runCommand("git", "fetch", "--all")
}

// Amend the last commit with the message
func GitCommitAmend(message string) error {
	return runCommand("git", "commit", "--amend", "-m", message)
}

// Amend the last commit but leave the message unchanged
func GitCommitAmendNoEdit() error {
	return runCommand("git", "commit", "--amend", "--no-edit")
}

// Add all the files to the index
func GitAddAll() error {
	return runCommand("git", "add", ".")
}

// Remove all the files from the index
func GitRemoveAll() error {
	return runCommand("git", "reset", "HEAD", "--", ".")
}

// Diff the index with the HEAD and return true if there is no diff
func GitDiffIndexHead() (bool, error) {
	// If there is no diff, output will be empty
	output, err := runCommandWithOutput("git", "diff", "--cached", "HEAD")
	if err != nil {
		return false, err
	}
	return output == "", runCommandWithStdin("git", "diff", "--cached", "HEAD")
}

//

// Diff the index with a given commit or branch and return true if there is no diff
func GitDiffIndexCommit(commit string) (bool, error) {
	// If there is no diff, output will be empty
	output, err := runCommandWithOutput("git", "diff", "--cached", commit)
	if err != nil {
		return false, err
	}
	return output == "", runCommandWithStdin("git", "diff", "--cached", commit)
}

// Diff a ref to another ref and return true if there is no diff
func GitDiffRef(ref1 string, ref2 string) (bool, error) {
	// If there is no diff, output will be empty
	output, err := runCommandWithOutput("git", "diff", ref1, ref2)
	if err != nil {
		return false, err
	}
	return output == "", runCommandWithStdin("git", "diff", ref1, ref2)

}

// Delete a remote branch
func GitDeleteRemoteBranch(remote string, branch string) error {
	return runCommand("git", "push", remote, "--delete", branch)
}

// Stash the changes
func GitStash() error {
	return runCommand("git", "stash")
}

// Pop the changes from the stash
func GitStashPop() error {
	return runCommand("git", "stash", "pop")
}
