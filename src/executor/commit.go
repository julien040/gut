package executor

import (
	"github.com/go-git/go-git/v5"
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
		if v.Worktree == git.Unmodified {
			continue
		}
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
