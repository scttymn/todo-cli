package pkg

import (
	"fmt"
	"os"
	"os/exec"

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

	// Use git command to create branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = wd
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s", branchName, string(output))
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

	// Check if branch exists
	branchRef := plumbing.NewBranchReferenceName(branchName)
	_, err = repo.Reference(branchRef, true)
	if err != nil {
		return fmt.Errorf("branch %s does not exist: %w", branchName, err)
	}

	// Use git command directly to avoid working directory changes
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = wd
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %s", branchName, string(output))
	}

	return nil
}