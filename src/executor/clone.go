package executor

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
)

func Clone(repo string, path string, oldCommitSave bool) error {
	options := &git.CloneOptions{
		URL: repo,
	}
	if !oldCommitSave {
		options.Depth = 1
	}
	_, err := git.PlainClone(path, false, options)
	if err != nil {
		return err
	}
	return nil
}

func CloneWithAuth(repo string, path string, username string, password string, oldCommitSave bool) error {
	options := &git.CloneOptions{
		URL: repo,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	}
	if !oldCommitSave {
		options.Depth = 1
	}
	_, err := git.PlainClone(path, false, options)
	if err != nil {
		return err
	}
	return nil
}
