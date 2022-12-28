package executor

import (
	"github.com/go-git/go-git/v5/config"
)

type Remote struct {
	Name string
	Url  string
}

func ListRemote(path string) ([]Remote, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}
	var remoteList []Remote
	for _, remote := range remotes {
		remoteList = append(remoteList, Remote{
			Name: remote.Config().Name,
			Url:  remote.Config().URLs[0],
		})
	}
	return remoteList, nil
}

func AddRemote(path string, name string, url string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})

	return err
}

func RemoveRemote(path string, name string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	err = repo.DeleteRemote(name)
	return err
}
