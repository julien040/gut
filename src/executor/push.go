package executor

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
)

func Push(path string, remote string, username string, password string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	err = repo.Push(&git.PushOptions{
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
