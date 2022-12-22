package check

import (
	"net/url"
)

// Check if repo URL is from GitHub
func IsGitHubRepo(repo string) (bool, error) {
	// Parse the repo URL
	parsedRepo, err := url.Parse(repo)
	if err != nil {
		return false, err
	}
	// Check if the host is github.com
	if parsedRepo.Host == "github.com" {
		return true, nil
	}
	return false, nil
}

type Status string

const (
	Public   Status = "available"
	Private  Status = "require_auth"
	NotFound Status = "not_found"
)

/* // Check if the repo is not found, requires authentication or is public
func checkRepoStatus(repo string) (Status, error) {
	// Check if the repo is from GitHub
	isGitHub, err := IsGitHubRepo(repo)
	if err != nil {
		return NotFound, err
	}
	if !isGitHub {
		return NotFound, nil
	}
	// Check if the repo is public

} */
