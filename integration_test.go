package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupIntegrationTest creates a temporary directory and builds the CLI binary
func setupIntegrationTest(t *testing.T) (string, string) {
	// Create temp directory
	testDir, err := os.MkdirTemp("", "todo-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Build the CLI binary
	binaryPath := filepath.Join(testDir, "todo")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
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

	// Initialize git repo for testing
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	
	// Create initial commit
	os.WriteFile("README.md", []byte("# Test repo"), 0644)
	exec.Command("git", "add", "README.md").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	// Setup cleanup
	t.Cleanup(func() {
		os.Chdir(originalDir)
		os.RemoveAll(testDir)
	})

	return testDir, binaryPath
}

// runCLI executes the CLI binary with given arguments and returns stdout, stderr, and exit code
func runCLI(t *testing.T, binaryPath string, args ...string) (string, string, int) {
	cmd := exec.Command(binaryPath, args...)
	
	stdout, err := cmd.Output()
	var stderr []byte
	if exitError, ok := err.(*exec.ExitError); ok {
		stderr = exitError.Stderr
		return string(stdout), string(stderr), exitError.ExitCode()
	} else if err != nil {
		t.Fatalf("Failed to run CLI command: %v", err)
	}
	
	return string(stdout), "", 0
}

// runCLIWithInput executes the CLI binary with stdin input
func runCLIWithInput(t *testing.T, binaryPath string, input string, args ...string) (string, string, int) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(input)
	
	stdout, err := cmd.Output()
	var stderr []byte
	if exitError, ok := err.(*exec.ExitError); ok {
		stderr = exitError.Stderr
		return string(stdout), string(stderr), exitError.ExitCode()
	} else if err != nil {
		t.Fatalf("Failed to run CLI command with input: %v", err)
	}
	
	return string(stdout), "", 0
}

func TestListCommand(t *testing.T) {
	_, binaryPath := setupIntegrationTest(t)
	
	// Test creating a new list
	stdout, stderr, exitCode := runCLIWithInput(t, binaryPath, "y\n", "list", "authentication")
	
	if exitCode != 0 {
		t.Fatalf("list command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	// Check that list was created and switched
	if !strings.Contains(stdout, "Created and switched to list 'authentication'") {
		t.Errorf("Expected list creation message, got: %s", stdout)
	}
	
	// Check that todo file was created
	if !strings.Contains(stdout, "Initialized todo file: .todo/authentication.md") {
		t.Errorf("Expected todo file creation message, got: %s", stdout)
	}
	
	// Verify the todo file exists
	todoFile := ".todo/authentication.md"
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		t.Error("Todo file was not created")
	}
	
	// Test switching to existing list
	stdout, stderr, exitCode = runCLIWithInput(t, binaryPath, "y\n", "list", "authentication")
	
	if exitCode != 0 {
		t.Fatalf("list command failed for existing list with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Switched to existing list 'authentication'") {
		t.Errorf("Expected list switch message, got: %s", stdout)
	}
	
	// Test showing all lists
	stdout, stderr, exitCode = runCLI(t, binaryPath, "list")
	
	if exitCode != 0 {
		t.Fatalf("list command failed to show all lists with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Lists:") {
		t.Errorf("Expected lists header, got: %s", stdout)
	}
}

func TestAddCheckUncheckWorkflow(t *testing.T) {
	_, binaryPath := setupIntegrationTest(t)
	
	// Create a list
	runCLIWithInput(t, binaryPath, "y\n", "list", "testing")
	
	// Add some todo items
	stdout, stderr, exitCode := runCLI(t, binaryPath, "add", "First todo item")
	if exitCode != 0 {
		t.Fatalf("add command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Added todo item to feature 'testing': First todo item") {
		t.Errorf("Expected add confirmation, got: %s", stdout)
	}
	
	runCLI(t, binaryPath, "add", "Second todo item")
	
	// Check progress
	stdout, stderr, exitCode = runCLI(t, binaryPath, "progress")
	if exitCode != 0 {
		t.Fatalf("progress command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Progress: 0/2 completed") {
		t.Errorf("Expected 0/2 progress, got: %s", stdout)
	}
	
	// Check off first item
	stdout, stderr, exitCode = runCLI(t, binaryPath, "check", "1")
	if exitCode != 0 {
		t.Fatalf("check command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Marked item 1 as completed") {
		t.Errorf("Expected check confirmation, got: %s", stdout)
	}
	
	// Check progress again
	stdout, stderr, exitCode = runCLI(t, binaryPath, "progress")
	if exitCode != 0 {
		t.Fatalf("progress command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Progress: 1/2 completed") {
		t.Errorf("Expected 1/2 progress, got: %s", stdout)
	}
	
	// Uncheck the item
	stdout, stderr, exitCode = runCLI(t, binaryPath, "uncheck", "1")
	if exitCode != 0 {
		t.Fatalf("uncheck command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Marked item 1 as not completed") {
		t.Errorf("Expected uncheck confirmation, got: %s", stdout)
	}
	
	// Check progress is back to 0/2
	stdout, stderr, exitCode = runCLI(t, binaryPath, "progress")
	if exitCode != 0 {
		t.Fatalf("progress command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Progress: 0/2 completed") {
		t.Errorf("Expected 0/2 progress after uncheck, got: %s", stdout)
	}
}

func TestProgressCommandWithAll(t *testing.T) {
	_, binaryPath := setupIntegrationTest(t)
	
	// Create multiple lists with different progress
	runCLIWithInput(t, binaryPath, "y\n", "list", "feature1")
	runCLI(t, binaryPath, "add", "Feature 1 todo")
	runCLI(t, binaryPath, "check", "1")
	
	runCLIWithInput(t, binaryPath, "y\n", "list", "feature2")
	runCLI(t, binaryPath, "add", "Feature 2 todo 1")
	runCLI(t, binaryPath, "add", "Feature 2 todo 2")
	// Leave both unchecked
	
	// Test progress --all
	stdout, stderr, exitCode := runCLI(t, binaryPath, "progress", "--all")
	if exitCode != 0 {
		t.Fatalf("progress --all command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	// Should show both features
	if !strings.Contains(stdout, "feature1 - 1/1 completed (100%)") {
		t.Errorf("Expected feature1 100%% complete, got: %s", stdout)
	}
	
	if !strings.Contains(stdout, "feature2 - 0/2 completed (0%)") {
		t.Errorf("Expected feature2 0%% complete, got: %s", stdout)
	}
	
	// Test short flag
	stdout, stderr, exitCode = runCLI(t, binaryPath, "progress", "-a")
	if exitCode != 0 {
		t.Fatalf("progress -a command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "Lists:") {
		t.Errorf("Expected lists header, got: %s", stdout)
	}
}

func TestVersionCommand(t *testing.T) {
	_, binaryPath := setupIntegrationTest(t)
	
	stdout, stderr, exitCode := runCLI(t, binaryPath, "version")
	if exitCode != 0 {
		t.Fatalf("version command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	if !strings.Contains(stdout, "todo CLI v0.1.0") {
		t.Errorf("Expected version string, got: %s", stdout)
	}
}

func TestHelpCommand(t *testing.T) {
	_, binaryPath := setupIntegrationTest(t)
	
	stdout, stderr, exitCode := runCLI(t, binaryPath, "--help")
	if exitCode != 0 {
		t.Fatalf("help command failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	
	// Check for key commands
	expectedCommands := []string{"list", "add", "check", "uncheck", "progress", "version"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("Expected to find command %s in help output, got: %s", cmd, stdout)
		}
	}
}