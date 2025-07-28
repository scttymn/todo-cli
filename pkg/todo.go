package pkg

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type TodoItem struct {
	ID            int
	Text          string
	Completed     bool
	CompletedTime *time.Time
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
	
	// Updated regex to capture optional timestamp: - [x] task text (completed: 2024-01-15 10:30)
	checkboxRegex := regexp.MustCompile(`^- \[([ x])\] (.+?)(?:\s+\(completed:\s+(.+?)\))?$`)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if match := checkboxRegex.FindStringSubmatch(line); match != nil {
			completed := match[1] == "x"
			text := match[2]
			var completedTime *time.Time
			
			// Parse timestamp if present
			if completed && len(match) > 3 && match[3] != "" {
				if parsedTime, err := time.Parse("2006-01-02 15:04", match[3]); err == nil {
					completedTime = &parsedTime
				}
			}
			
			items = append(items, TodoItem{
				ID:            itemID,
				Text:          text,
				Completed:     completed,
				CompletedTime: completedTime,
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
			if item.CompletedTime != nil {
				fmt.Fprintf(file, "- [%s] %s (completed: %s)\n", checkbox, item.Text, item.CompletedTime.Format("2006-01-02 15:04"))
			} else {
				fmt.Fprintf(file, "- [%s] %s\n", checkbox, item.Text)
			}
		} else {
			fmt.Fprintf(file, "- [%s] %s\n", checkbox, item.Text)
		}
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
		ID:            newID,
		Text:          text,
		Completed:     false,
		CompletedTime: nil,
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

	now := time.Now()
	todoList.Items[itemID-1].Completed = true
	todoList.Items[itemID-1].CompletedTime = &now
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
	todoList.Items[itemID-1].CompletedTime = nil
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

	fmt.Println("Lists:")
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

func ShowHistory() error {
	if err := EnsureTodoDirectory(); err != nil {
		return fmt.Errorf("failed to ensure .todo directory: %w", err)
	}

	files, err := os.ReadDir(".todo")
	if err != nil {
		return fmt.Errorf("failed to read .todo directory: %w", err)
	}

	type CompletedItem struct {
		Text      string
		List      string
		Completed time.Time
	}

	var completedItems []CompletedItem

	// Collect all completed items from all lists
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			listName := strings.TrimSuffix(file.Name(), ".md")
			
			todoList, err := ParseTodoFile(listName)
			if err != nil {
				continue // Skip files we can't parse
			}

			for _, item := range todoList.Items {
				if item.Completed && item.CompletedTime != nil {
					completedItems = append(completedItems, CompletedItem{
						Text:      item.Text,
						List:      listName,
						Completed: *item.CompletedTime,
					})
				}
			}
		}
	}

	if len(completedItems) == 0 {
		fmt.Println("No completed todos found.")
		return nil
	}

	// Sort by completion time (newest first)
	for i := 0; i < len(completedItems); i++ {
		for j := i + 1; j < len(completedItems); j++ {
			if completedItems[i].Completed.Before(completedItems[j].Completed) {
				completedItems[i], completedItems[j] = completedItems[j], completedItems[i]
			}
		}
	}

	fmt.Println("Completed Todo History:")
	fmt.Println()

	currentDate := ""
	for _, item := range completedItems {
		itemDate := item.Completed.Format("2006-01-02")
		if itemDate != currentDate {
			if currentDate != "" {
				fmt.Println()
			}
			fmt.Printf("📅 %s\n", item.Completed.Format("Monday, January 2, 2006"))
			currentDate = itemDate
		}
		
		timeStr := item.Completed.Format("15:04")
		fmt.Printf("  ✅ %s [%s] (%s)\n", item.Text, item.List, timeStr)
	}

	return nil
}

func EditTodoFile(branchName string) error {
	// Get the editor from environment variable
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable is not set. Please set it to your preferred editor (e.g., export EDITOR=nvim)")
	}
	
	// Ensure the todo file exists
	if !TodoFileExists(branchName) {
		err := CreateTodoFile(branchName)
		if err != nil {
			return fmt.Errorf("failed to create todo file: %w", err)
		}
	}
	
	// Get the file path
	filePath := GetTodoFilePath(branchName)
	
	// Execute the editor command
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run editor %s: %w", editor, err)
	}
	
	return nil
}