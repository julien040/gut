package executor

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SoftReset(path string, commit object.Commit) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Reset(&git.ResetOptions{
		Mode:   git.SoftReset,
		Commit: commit.Hash,
	})
	return err

}
