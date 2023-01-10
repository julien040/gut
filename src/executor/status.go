package executor

type FileStatus struct {
	Path   string
	Status string
}

func GetStatus(path string) ([]FileStatus, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	AddAll(path)
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	// Add all files without one listed in gitignore

	status, err := w.Status()
	if err != nil {
		return nil, err
	}
	var fileStatuses []FileStatus
	for file, status := range status {
		fileStatuses = append(fileStatuses, FileStatus{Path: file, Status: string(status.Staging)})
	}
	return fileStatuses, nil

}
