package executor

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
)

func Clone(repo string, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repo,
	})
	if err != nil {
		return err
	}
	return nil
}

func CloneWithAuth(repo string, path string, username string, password string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repo,
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
