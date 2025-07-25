package pkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetCurrentBranch() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		return "", fmt.Errorf("not a git repository or unable to open: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	return head.Hash().String()[:7], nil
}

func CreateBranch(branchName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		return fmt.Errorf("not a git repository or unable to open: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(branchRef, head.Hash())

	err = repo.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func SwitchBranch(branchName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		return fmt.Errorf("not a git repository or unable to open: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
	})
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %w", branchName, err)
	}

	return nil
}