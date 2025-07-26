package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		return "", fmt.Errorf("this directory is not a git repository. Please run 'git init' or navigate to a git repository")
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

func GetFeatureName() (string, error) {
	branchName, err := GetCurrentBranch()
	if err != nil {
		return "", err
	}
	
	// If it's a feature branch (feature/name), extract just the feature name
	if strings.HasPrefix(branchName, "feature/") {
		return strings.TrimPrefix(branchName, "feature/"), nil
	}
	
	// Otherwise return the full branch name
	return branchName, nil
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
		if strings.Contains(string(output), "not a git repository") {
			return fmt.Errorf("this directory is not a git repository. Please run 'git init' first")
		}
		return fmt.Errorf("failed to create branch %s: %s", branchName, string(output))
	}

	return nil
}

func BranchExists(branchName string) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		return false, fmt.Errorf("this directory is not a git repository. Please run 'git init' or navigate to a git repository")
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	_, err = repo.Reference(branchRef, true)
	if err != nil {
		return false, nil // Branch doesn't exist
	}

	return true, nil
}

func HasUncommittedChanges() (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Use git status --porcelain to check for changes
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = wd
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("unable to check git status. Make sure you're in a git repository")
	}

	// If output is not empty, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func SwitchBranch(branchName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		return fmt.Errorf("this directory is not a git repository. Please run 'git init' or navigate to a git repository")
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
		if strings.Contains(string(output), "not a git repository") {
			return fmt.Errorf("this directory is not a git repository. Please run 'git init' first")
		}
		return fmt.Errorf("failed to switch to branch %s: %s", branchName, string(output))
	}

	return nil
}

func DeleteBranch(branchName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Use git command to delete branch
	cmd := exec.Command("git", "branch", "-D", branchName)
	cmd.Dir = wd
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "not a git repository") {
			return fmt.Errorf("this directory is not a git repository. Please run 'git init' first")
		}
		return fmt.Errorf("failed to delete branch %s: %s", branchName, string(output))
	}

	return nil
}