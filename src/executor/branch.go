package executor

import (
	"errors"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func ListBranches(path string) ([]string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	branchesIter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	var branches []*plumbing.Reference

	err = branchesIter.ForEach(func(branch *plumbing.Reference) error {
		branches = append(branches, branch)
		return nil
	})

	// Sort by commit time
	sort.Slice(branches, func(i, j int) bool {
		commitI, err := repo.CommitObject(branches[i].Hash())
		if err != nil {
			return false
		}
		commitJ, err := repo.CommitObject(branches[j].Hash())
		if err != nil {
			return false
		}
		return commitI.Committer.When.After(commitJ.Committer.When)
	})

	var branchNames []string

	for _, branch := range branches {
		branchNames = append(branchNames, branch.Name().Short())
	}
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
		Keep:   true,
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
		Force:  true,
	})
	return err

}

func DeleteBranch(path string, branchName string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}
	// Check if branch exists
	exists, err := CheckIfBranchExists(path, branchName)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("branch does not exist")
	}

	// Check if branch is current branch
	currentBranch, err := GetCurrentBranch(path)
	if err != nil {
		return err
	}
	if currentBranch == branchName {
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
	// List all branches and check if branchName is in the list
	branches, err := repo.Branches()
	if err != nil {
		return false, err

	}
	var branchNames []string

	err = branches.ForEach(func(branch *plumbing.Reference) error {
		branchNames = append(branchNames, branch.Name().Short())
		return nil
	})
	if err != nil {
		return false, err
	}
	for _, branch := range branchNames {
		if branch == branchName {
			return true, nil
		}
	}
	return false, nil

}

func GetCurrentBranch(path string) (string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return "", err
	}
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}
