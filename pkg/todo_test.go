package pkg

import (
	"os"
	"testing"
)

func setupTestDir(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "todo-test-*")
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
	
	// Store original directory for cleanup
	t.Cleanup(func() {
		os.Chdir(originalDir)
		os.RemoveAll(testDir)
	})
	
	return testDir
}

func TestGetTodoFilePath(t *testing.T) {
	tests := []struct {
		branchName string
		expected   string
	}{
		{"authentication", ".todo/authentication.md"},
		{"payment-system", ".todo/payment-system.md"},
		{"main", ".todo/main.md"},
	}
	
	for _, tt := range tests {
		t.Run(tt.branchName, func(t *testing.T) {
			result := GetTodoFilePath(tt.branchName)
			if result != tt.expected {
				t.Errorf("GetTodoFilePath(%q) = %q, want %q", tt.branchName, result, tt.expected)
			}
		})
	}
}

func TestTodoFileExists(t *testing.T) {
	setupTestDir(t)
	
	// Test non-existent file
	if TodoFileExists("nonexistent") {
		t.Error("TodoFileExists should return false for non-existent file")
	}
	
	// Create .todo directory and file
	err := EnsureTodoDirectory()
	if err != nil {
		t.Fatalf("Failed to create .todo directory: %v", err)
	}
	
	err = CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	
	// Test existing file
	if !TodoFileExists("test-feature") {
		t.Error("TodoFileExists should return true for existing file")
	}
}

func TestCreateTodoFile(t *testing.T) {
	setupTestDir(t)
	
	err := CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("CreateTodoFile failed: %v", err)
	}
	
	// Check if file exists
	filePath := GetTodoFilePath("test-feature")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Todo file was not created")
	}
	
	// Check file contents
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read todo file: %v", err)
	}
	
	expected := "# Todo List for test-feature\n\n"
	if string(content) != expected {
		t.Errorf("File content = %q, want %q", string(content), expected)
	}
}

func TestParseTodoFile(t *testing.T) {
	setupTestDir(t)
	
	// Create a test todo file with some items
	err := EnsureTodoDirectory()
	if err != nil {
		t.Fatalf("Failed to create .todo directory: %v", err)
	}
	
	testContent := `# Todo List for test-feature

- [ ] First todo item
- [x] Completed todo item
- [ ] Another pending item
`
	
	filePath := GetTodoFilePath("test-feature")
	err = os.WriteFile(filePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Parse the file
	todoList, err := ParseTodoFile("test-feature")
	if err != nil {
		t.Fatalf("ParseTodoFile failed: %v", err)
	}
	
	// Check parsed items
	expected := []TodoItem{
		{ID: 1, Text: "First todo item", Completed: false},
		{ID: 2, Text: "Completed todo item", Completed: true},
		{ID: 3, Text: "Another pending item", Completed: false},
	}
	
	if len(todoList.Items) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(todoList.Items))
	}
	
	for i, item := range todoList.Items {
		if item.ID != expected[i].ID {
			t.Errorf("Item %d: ID = %d, want %d", i, item.ID, expected[i].ID)
		}
		if item.Text != expected[i].Text {
			t.Errorf("Item %d: Text = %q, want %q", i, item.Text, expected[i].Text)
		}
		if item.Completed != expected[i].Completed {
			t.Errorf("Item %d: Completed = %v, want %v", i, item.Completed, expected[i].Completed)
		}
	}
}

func TestAddTodoItem(t *testing.T) {
	setupTestDir(t)
	
	err := CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	
	// Add a todo item
	err = AddTodoItem("test-feature", "Test todo item")
	if err != nil {
		t.Fatalf("AddTodoItem failed: %v", err)
	}
	
	// Parse and verify
	todoList, err := ParseTodoFile("test-feature")
	if err != nil {
		t.Fatalf("ParseTodoFile failed: %v", err)
	}
	
	if len(todoList.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(todoList.Items))
	}
	
	item := todoList.Items[0]
	if item.ID != 1 {
		t.Errorf("ID = %d, want 1", item.ID)
	}
	if item.Text != "Test todo item" {
		t.Errorf("Text = %q, want %q", item.Text, "Test todo item")
	}
	if item.Completed {
		t.Error("New item should not be completed")
	}
}

func TestCheckTodoItem(t *testing.T) {
	setupTestDir(t)
	
	err := CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	
	// Add some items
	err = AddTodoItem("test-feature", "First item")
	if err != nil {
		t.Fatalf("AddTodoItem failed: %v", err)
	}
	err = AddTodoItem("test-feature", "Second item")
	if err != nil {
		t.Fatalf("AddTodoItem failed: %v", err)
	}
	
	// Check the first item
	err = CheckTodoItem("test-feature", 1)
	if err != nil {
		t.Fatalf("CheckTodoItem failed: %v", err)
	}
	
	// Verify the item is checked
	todoList, err := ParseTodoFile("test-feature")
	if err != nil {
		t.Fatalf("ParseTodoFile failed: %v", err)
	}
	
	if !todoList.Items[0].Completed {
		t.Error("First item should be completed")
	}
	if todoList.Items[1].Completed {
		t.Error("Second item should not be completed")
	}
}

func TestUncheckTodoItem(t *testing.T) {
	setupTestDir(t)
	
	err := CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	
	// Add and check an item
	err = AddTodoItem("test-feature", "Test item")
	if err != nil {
		t.Fatalf("AddTodoItem failed: %v", err)
	}
	err = CheckTodoItem("test-feature", 1)
	if err != nil {
		t.Fatalf("CheckTodoItem failed: %v", err)
	}
	
	// Uncheck the item
	err = UncheckTodoItem("test-feature", 1)
	if err != nil {
		t.Fatalf("UncheckTodoItem failed: %v", err)
	}
	
	// Verify the item is unchecked
	todoList, err := ParseTodoFile("test-feature")
	if err != nil {
		t.Fatalf("ParseTodoFile failed: %v", err)
	}
	
	if todoList.Items[0].Completed {
		t.Error("Item should not be completed after unchecking")
	}
}

func TestCheckTodoItemInvalidID(t *testing.T) {
	setupTestDir(t)
	
	err := CreateTodoFile("test-feature")
	if err != nil {
		t.Fatalf("Failed to create todo file: %v", err)
	}
	
	// Try to check non-existent item
	err = CheckTodoItem("test-feature", 999)
	if err == nil {
		t.Error("CheckTodoItem should fail for invalid ID")
	}
	
	// Try to check item with ID 0
	err = CheckTodoItem("test-feature", 0)
	if err == nil {
		t.Error("CheckTodoItem should fail for ID 0")
	}
}