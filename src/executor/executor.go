package executor

import (
	"github.com/go-git/go-git/v5"
)

func IsPathGitRepo(path string) bool {
	_, err := git.PlainOpen(path)

	return err == nil
}

func OpenRepo(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func SetUserConfig(path string, name string, email string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	config, err := repo.Config()
	if err != nil {
		return err
	}
	config.User.Name = name
	config.User.Email = email
	err = repo.SetConfig(config)
	return err
}

func GetGitURL(path string) (string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return "", err
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	return remote.Config().URLs[0], nil
}
