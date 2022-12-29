package executor

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func ListBranches(path string) ([]string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}
	var branchNames []string
	err = branches.ForEach(func(branch *plumbing.Reference) error {
		branchNames = append(branchNames, branch.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return branchNames, nil
}

func CreateBranch(path string, branchName string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	// Check if branch exists
	_, err = repo.Branch(branchName)
	if err == nil {
		return errors.New("branch already exists")
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	return err
}

func CheckoutBranch(path string, branchName string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	return err

}

func DeleteBranch(path string, branchName string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	// Check if branch exists
	_, err = repo.Branch(branchName)
	if err != nil {
		return err
	}
	// Check if branch is current branch
	currentBranch, err := repo.Head()
	if err != nil {
		return err
	}
	if currentBranch.Name().Short() == branchName {
		return errors.New("cannot delete current branch")
	}
	branch := plumbing.NewBranchReferenceName(branchName)
	err = repo.Storer.RemoveReference(branch)
	return err
}

func CheckIfBranchExists(path string, branchName string) (bool, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return false, err
	}
	_, err = repo.Branch(branchName)
	if err != nil {
		return false, nil
	}
	return true, nil
}
