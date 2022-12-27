package executor

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type CommitResult struct {
	Hash         string
	FilesAdded   int
	FilesUpdated int
	FilesDeleted int
}

func Commit(path string, message string) (CommitResult, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return CommitResult{}, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return CommitResult{}, err
	}
	// Replace the .gitignore file with the one in the repo
	replaceGitIgnore(w, filepath.Join(path, ".gitignore"))

	// Add all files
	err = w.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return CommitResult{}, err
	}

	status, err := w.Status()
	if err != nil {
		return CommitResult{}, err
	}
	fileAdded := 0
	fileUpdated := 0
	fileDeleted := 0
	for _, v := range status {
		if v.Staging == git.Added {
			fileAdded++
		}
		if v.Staging == git.Deleted {
			fileDeleted++
		}
		if v.Staging == git.Modified {
			fileUpdated++
		}
	}
	hash, err := w.Commit(message, &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Hash:         hash.String(),
		FilesAdded:   fileAdded,
		FilesUpdated: fileUpdated,
		FilesDeleted: fileDeleted,
	}, nil
}

// Load the .gitignore file, and add to the "exclude" section of the worktree
//
// This is a workaround for the fact that go-git does not support .gitignore
func replaceGitIgnore(wk *git.Worktree, path string) {
	f, _ := os.Open(path)
	defer f.Close()

	// Read the file
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		// Add the line to the exclude list
		wk.Excludes = append(wk.Excludes, gitignore.ParsePattern(line, nil))
	}

}
