package pkg

import (
	"os"
	"os/exec"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "todo-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	
	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = testDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	// Configure git (needed for commits)
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	
	// Create initial commit
	err = os.WriteFile("README.md", []byte("# Test repo"), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}
	
	exec.Command("git", "add", "README.md").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()
	
	// Store original directory for cleanup
	t.Cleanup(func() {
		os.Chdir(originalDir)
		os.RemoveAll(testDir)
	})
	
	return testDir
}

func TestGetCurrentBranch(t *testing.T) {
	setupGitRepo(t)
	
	branch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	
	// Default branch should be main or master
	if branch != "main" && branch != "master" {
		t.Errorf("Expected branch 'main' or 'master', got %q", branch)
	}
}

func TestGetFeatureName(t *testing.T) {
	setupGitRepo(t)
	
	// Test current branch (should be main or master)
	t.Run("current branch", func(t *testing.T) {
		featureName, err := GetFeatureName()
		if err != nil {
			t.Fatalf("GetFeatureName failed: %v", err)
		}
		
		// Should return the actual branch name (main or master)
		if featureName != "main" && featureName != "master" {
			t.Errorf("GetFeatureName() = %q, want 'main' or 'master'", featureName)
		}
	})
	
	// Test feature branch extraction
	t.Run("feature branch", func(t *testing.T) {
		// Create a feature branch
		cmd := exec.Command("git", "checkout", "-b", "feature/authentication")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to create feature branch: %v", err)
		}
		
		featureName, err := GetFeatureName()
		if err != nil {
			t.Fatalf("GetFeatureName failed: %v", err)
		}
		
		if featureName != "authentication" {
			t.Errorf("GetFeatureName() = %q, want %q", featureName, "authentication")
		}
	})
}

func TestBranchExists(t *testing.T) {
	setupGitRepo(t)
	
	// Test existing branch
	exists, err := BranchExists("main")
	if err != nil {
		// Try master if main doesn't exist
		exists, err = BranchExists("master")
		if err != nil {
			t.Fatalf("BranchExists failed: %v", err)
		}
	}
	
	if !exists {
		t.Error("BranchExists should return true for main/master branch")
	}
	
	// Test non-existent branch
	exists, err = BranchExists("nonexistent-branch")
	if err != nil {
		t.Fatalf("BranchExists failed for non-existent branch: %v", err)
	}
	
	if exists {
		t.Error("BranchExists should return false for non-existent branch")
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	setupGitRepo(t)
	
	// Test clean repository
	hasChanges, err := HasUncommittedChanges()
	if err != nil {
		t.Fatalf("HasUncommittedChanges failed: %v", err)
	}
	
	if hasChanges {
		t.Error("HasUncommittedChanges should return false for clean repo")
	}
	
	// Create uncommitted change
	err = os.WriteFile("test.txt", []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test dirty repository
	hasChanges, err = HasUncommittedChanges()
	if err != nil {
		t.Fatalf("HasUncommittedChanges failed: %v", err)
	}
	
	if !hasChanges {
		t.Error("HasUncommittedChanges should return true for dirty repo")
	}
}

func TestCreateBranch(t *testing.T) {
	setupGitRepo(t)
	
	// Create a new branch
	err := CreateBranch("feature/test-branch")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}
	
	// Verify we're on the new branch
	branch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	
	if branch != "feature/test-branch" {
		t.Errorf("Expected to be on 'feature/test-branch', got %q", branch)
	}
	
	// Verify branch exists
	exists, err := BranchExists("feature/test-branch")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	
	if !exists {
		t.Error("Branch should exist after creation")
	}
}

func TestSwitchBranch(t *testing.T) {
	setupGitRepo(t)
	
	// Get the current (main) branch name
	mainBranch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	
	// Create a branch to switch to
	err = CreateBranch("feature/switch-test")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}
	
	// Switch back to main
	err = SwitchBranch(mainBranch)
	if err != nil {
		t.Fatalf("SwitchBranch to main failed: %v", err)
	}
	
	// Verify we're on main
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	
	if currentBranch != mainBranch {
		t.Errorf("Expected to be on %q, got %q", mainBranch, currentBranch)
	}
	
	// Switch back to feature branch
	err = SwitchBranch("feature/switch-test")
	if err != nil {
		t.Fatalf("SwitchBranch to feature failed: %v", err)
	}
	
	// Verify we're on the feature branch
	currentBranch, err = GetCurrentBranch()
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	
	if currentBranch != "feature/switch-test" {
		t.Errorf("Expected to be on 'feature/switch-test', got %q", currentBranch)
	}
}