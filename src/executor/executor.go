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

func GetUserConfig(path string) (string, string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return "", "", err
	}
	config, err := repo.Config()
	if err != nil {
		return "", "", err
	}
	return config.User.Name, config.User.Email, nil
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

func IsWorkTreeClean(path string) (bool, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return false, err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	AddAll(path)
	status, err := worktree.Status()
	if err != nil {
		return false, err
	}

	// See /src/executor/status.go:29
	for _, statusFile := range status {
		if statusFile.Staging != git.Untracked {
			return false, nil
		}
	}
	return true, nil
}

func IsDetachedHead(path string) (bool, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return false, err
	}
	head, err := repo.Head()
	if err != nil {
		return false, err
	}
	return head.Name() == "HEAD", nil
}
