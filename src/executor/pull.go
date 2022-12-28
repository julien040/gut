package executor

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
)

func Pull(path string, remote string, username string, password string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		RemoteName: remote,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		return err
	}
	return nil

}
