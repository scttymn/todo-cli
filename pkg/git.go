package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func IsGitRepository() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}
	
	_, err = git.PlainOpen(wd)
	return err == nil
}

func HasCommits() bool {
	if !IsGitRepository() {
		return false
	}
	
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) != "0"
}

func InitTodoRepository() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize git repository if not already one
	if !IsGitRepository() {
		cmd := exec.Command("git", "init")
		cmd.Dir = wd
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %s", string(output))
		}
	}

	// Create .gitignore if it doesn't exist
	gitignorePath := ".gitignore"
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		gitignoreContent := `# Todo directory (local only)
.todo/

# Go build artifacts
todo-cli
todo
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work
`
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}
	}

	// Create initial README if it doesn't exist
	readmePath := "README.md"
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readmeContent := "# Project\n\nTodo lists are managed with the `todo` CLI tool.\n"
		if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
			return fmt.Errorf("failed to create README.md: %w", err)
		}
	}

	// Add files and make initial commit if no commits exist
	if !HasCommits() {
		// Add files
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = wd
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add files: %s", string(output))
		}

		// Make initial commit
		cmd = exec.Command("git", "commit", "-m", "Initial commit")
		cmd.Dir = wd
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to make initial commit: %s", string(output))
		}
	}

	return nil
}

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