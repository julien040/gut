package executor

import (
	"github.com/go-git/go-git/v5"
)

type FileStatus struct {
	Path   string
	Status string
}

func GetStatus(path string) ([]FileStatus, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	// Run git add . to add all files using git
	// or use go-git to add all files
	AddAll(path)
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := w.Status()
	if err != nil {
		return nil, err
	}
	/**
	 * Due to a bug in go-git, the status of untracked files is not returned.
	 * Sometimes, go-git returns an untracked status while git returns ignores the files
	 * This is a workaround to fix this issue.
	 * I thinks it's caused by the fact that go-git does not support nested .gitignore properly
	 */

	var fileStatuses []FileStatus
	for file, statusFile := range status {
		if statusFile.Staging == git.Untracked {
			continue
		}
		fileStatuses = append(fileStatuses, FileStatus{Path: file, Status: string(statusFile.Staging)})
	}
	return fileStatuses, nil

}
