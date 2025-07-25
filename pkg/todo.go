package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TodoItem struct {
	ID        int
	Text      string
	Completed bool
}

type TodoList struct {
	Items []TodoItem
}

func GetTodoFilePath(branchName string) string {
	return filepath.Join(".todo", branchName+".md")
}

func TodoFileExists(featureName string) bool {
	filePath := GetTodoFilePath(featureName)
	_, err := os.Stat(filePath)
	return err == nil
}

func EnsureTodoDirectory() error {
	return os.MkdirAll(".todo", 0755)
}

func CreateTodoFile(branchName string) error {
	if err := EnsureTodoDirectory(); err != nil {
		return fmt.Errorf("failed to create .todo directory: %w", err)
	}

	filePath := GetTodoFilePath(branchName)
	
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create todo file: %w", err)
	}
	defer file.Close()

	content := fmt.Sprintf("# Todo List for %s\n\n", branchName)
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write initial content: %w", err)
	}

	return nil
}

func ParseTodoFile(branchName string) (*TodoList, error) {
	filePath := GetTodoFilePath(branchName)
	
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &TodoList{Items: []TodoItem{}}, nil
		}
		return nil, fmt.Errorf("failed to open todo file: %w", err)
	}
	defer file.Close()

	var items []TodoItem
	scanner := bufio.NewScanner(file)
	itemID := 1
	
	checkboxRegex := regexp.MustCompile(`^- \[([ x])\] (.+)$`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if match := checkboxRegex.FindStringSubmatch(line); match != nil {
			completed := match[1] == "x"
			text := match[2]
			
			items = append(items, TodoItem{
				ID:        itemID,
				Text:      text,
				Completed: completed,
			})
			itemID++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading todo file: %w", err)
	}

	return &TodoList{Items: items}, nil
}

func WriteTodoFile(branchName string, todoList *TodoList) error {
	if err := EnsureTodoDirectory(); err != nil {
		return fmt.Errorf("failed to create .todo directory: %w", err)
	}

	filePath := GetTodoFilePath(branchName)
	
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create todo file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "# Todo List for %s\n\n", branchName)
	
	for _, item := range todoList.Items {
		checkbox := " "
		if item.Completed {
			checkbox = "x"
		}
		fmt.Fprintf(file, "- [%s] %s\n", checkbox, item.Text)
	}

	return nil
}

func AddTodoItem(branchName, text string) error {
	todoList, err := ParseTodoFile(branchName)
	if err != nil {
		return fmt.Errorf("failed to parse todo file: %w", err)
	}

	newID := len(todoList.Items) + 1
	todoList.Items = append(todoList.Items, TodoItem{
		ID:        newID,
		Text:      text,
		Completed: false,
	})

	return WriteTodoFile(branchName, todoList)
}

func CheckTodoItem(branchName string, itemID int) error {
	todoList, err := ParseTodoFile(branchName)
	if err != nil {
		return fmt.Errorf("failed to parse todo file: %w", err)
	}

	if itemID < 1 || itemID > len(todoList.Items) {
		return fmt.Errorf("invalid item ID: %d", itemID)
	}

	todoList.Items[itemID-1].Completed = true
	return WriteTodoFile(branchName, todoList)
}

func UncheckTodoItem(branchName string, itemID int) error {
	todoList, err := ParseTodoFile(branchName)
	if err != nil {
		return fmt.Errorf("failed to parse todo file: %w", err)
	}

	if itemID < 1 || itemID > len(todoList.Items) {
		return fmt.Errorf("invalid item ID: %d", itemID)
	}

	todoList.Items[itemID-1].Completed = false
	return WriteTodoFile(branchName, todoList)
}

func DisplayTodoList(branchName string) error {
	todoList, err := ParseTodoFile(branchName)
	if err != nil {
		return fmt.Errorf("failed to parse todo file: %w", err)
	}

	if len(todoList.Items) == 0 {
		fmt.Printf("No todos for branch '%s'\n", branchName)
		return nil
	}

	fmt.Printf("Todo list for branch '%s':\n\n", branchName)
	
	completed := 0
	for _, item := range todoList.Items {
		status := "[ ]"
		if item.Completed {
			status = "[x]"
			completed++
		}
		fmt.Printf("%d. %s %s\n", item.ID, status, item.Text)
	}

	fmt.Printf("\nProgress: %d/%d completed\n", completed, len(todoList.Items))
	return nil
}

func ListAllFeatures() error {
	if err := EnsureTodoDirectory(); err != nil {
		return fmt.Errorf("failed to ensure .todo directory: %w", err)
	}

	files, err := os.ReadDir(".todo")
	if err != nil {
		return fmt.Errorf("failed to read .todo directory: %w", err)
	}

	var features []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			featureName := strings.TrimSuffix(file.Name(), ".md")
			features = append(features, featureName)
		}
	}

	if len(features) == 0 {
		fmt.Println("No features found")
		return nil
	}

	fmt.Println("Features:")
	fmt.Println()

	for _, feature := range features {
		todoList, err := ParseTodoFile(feature)
		if err != nil {
			fmt.Printf("  %s - Error reading file: %v\n", feature, err)
			continue
		}

		completed := 0
		for _, item := range todoList.Items {
			if item.Completed {
				completed++
			}
		}

		total := len(todoList.Items)
		if total == 0 {
			fmt.Printf("  %s - No todos\n", feature)
		} else {
			percentage := (completed * 100) / total
			fmt.Printf("  %s - %d/%d completed (%d%%)\n", feature, completed, total, percentage)
		}
	}

	return nil
}