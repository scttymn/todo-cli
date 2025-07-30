package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scttymn/todo-cli/pkg"
	"github.com/spf13/cobra"
)

func requiresInit() bool {
	// Just ensure .todo directory exists
	if err := pkg.EnsureTodoDirectory(); err != nil {
		fmt.Printf("Failed to create .todo directory: %v\n", err)
		return true
	}
	return false
}

var rootCmd = &cobra.Command{
	Use:   "todo [command] [flags]",
	Short: "A CLI tool for managing todo lists",
	Long:  `todo is a CLI tool that manages todo lists in markdown files, helping you track tasks for different projects or features.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize todo management in the current directory",
	Long:  `Initialize the current directory for todo management by creating the .todo directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.EnsureTodoDirectory()
		if err != nil {
			fmt.Printf("Failed to initialize todo directory: %v\n", err)
			return
		}
		
		fmt.Println("✅ Todo management initialized successfully!")
		fmt.Println("You can now create todo lists with: todo list <name>")
	},
}


var addCmd = &cobra.Command{
	Use:   "add [todo-item]",
	Short: "Add a todo item to the current list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		todoItem := args[0]
		
		currentList, err := pkg.GetCurrentList()
		if err != nil {
			fmt.Printf("Error getting current list: %v\n", err)
			return
		}
		
		err = pkg.AddTodoItem(currentList, todoItem)
		if err != nil {
			fmt.Printf("Error adding todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Added todo item to list '%s': %s\n", currentList, todoItem)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [item-number]",
	Short: "Mark a todo item as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		itemNumber := args[0]
		
		currentList, err := pkg.GetCurrentList()
		if err != nil {
			fmt.Printf("Error getting current list: %v\n", err)
			return
		}
		
		itemID, err := strconv.Atoi(itemNumber)
		if err != nil {
			fmt.Printf("Invalid item number: %s\n", itemNumber)
			return
		}
		
		err = pkg.CheckTodoItem(currentList, itemID)
		if err != nil {
			fmt.Printf("Error checking todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Marked item %d as completed in list '%s'\n", itemID, currentList)
	},
}

var uncheckCmd = &cobra.Command{
	Use:   "uncheck [item-number]",
	Short: "Mark a todo item as not completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		itemNumber := args[0]
		
		currentList, err := pkg.GetCurrentList()
		if err != nil {
			fmt.Printf("Error getting current list: %v\n", err)
			return
		}
		
		itemID, err := strconv.Atoi(itemNumber)
		if err != nil {
			fmt.Printf("Invalid item number: %s\n", itemNumber)
			return
		}
		
		err = pkg.UncheckTodoItem(currentList, itemID)
		if err != nil {
			fmt.Printf("Error unchecking todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Marked item %d as not completed in list '%s'\n", itemID, currentList)
	},
}

var progressCmd = &cobra.Command{
	Use:   "progress [list-name]",
	Short: "Show progress for current list, specific list, or all lists\n                Available flags: --all",
	Long:  `Show todo progress:\n\n  todo progress             Current list progress\n  todo progress <name>      Specific list progress\n  todo progress --all       All lists progress`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		showAll, _ := cmd.Flags().GetBool("all")
		
		if showAll {
			if len(args) > 0 {
				fmt.Println("Error: Cannot use --all flag with list name")
				return
			}
			err := pkg.ListAllFeatures()
			if err != nil {
				fmt.Printf("Error showing progress: %v\n", err)
				return
			}
		} else if len(args) == 1 {
			// Show progress for specific list
			listName := args[0]
			
			// Check if the list exists by checking if todo file exists
			if !pkg.TodoFileExists(listName) {
				fmt.Printf("List '%s' does not exist\n", listName)
				return
			}
			
			err := pkg.DisplayTodoList(listName)
			if err != nil {
				fmt.Printf("Error displaying todo list: %v\n", err)
				return
			}
		} else {
			// Show progress for current list
			currentList, err := pkg.GetCurrentList()
			if err != nil {
				fmt.Printf("Error getting current list: %v\n", err)
				return
			}
			
			err = pkg.DisplayTodoList(currentList)
			if err != nil {
				fmt.Printf("Error displaying todo list: %v\n", err)
				return
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list [list-name]",
	Short: "Show all lists, switch to lists, or create new lists\n                Available flags: --delete",
	Long:  `Manage todo lists:\n\n  todo list                 Show all lists with progress\n  todo list <name>          Switch to or create list\n  todo list --delete <name> Delete list (requires confirmation)`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		deleteFlag, _ := cmd.Flags().GetBool("delete")
		
		if deleteFlag {
			if len(args) == 0 {
				fmt.Println("Error: --delete requires a list name")
				return
			}
			
			listName := args[0]
			
			// Check if we're currently on the list we're trying to delete
			currentList, err := pkg.GetCurrentList()
			if err != nil {
				fmt.Printf("Error getting current list: %v\n", err)
				return
			}
			
			if currentList == listName {
				fmt.Printf("Error: Cannot delete list '%s' because it is currently active.\n", listName)
				fmt.Println("Switch to another list first (e.g., 'todo list main')")
				return
			}
			
			// Check if list exists
			if !pkg.ListExists(listName) {
				fmt.Printf("List '%s' does not exist\n", listName)
				return
			}
			
			// Confirmation prompt
			fmt.Printf("Are you sure you want to delete list '%s'? This will remove the todo file. (y/N): ", listName)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return
			}
			
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Delete cancelled.")
				return
			}
			
			// Delete the todo file
			err = pkg.DeleteList(listName)
			if err != nil {
				fmt.Printf("Error deleting list: %v\n", err)
				return
			}
			
			fmt.Printf("Successfully deleted list '%s'\n", listName)
			return
		}
		
		if len(args) == 0 {
			// Show all lists
			err := pkg.ListAllFeatures()
			if err != nil {
				fmt.Printf("Error showing lists: %v\n", err)
				return
			}
		} else {
			// Switch to or create specific list
			listName := args[0]
			
			// Set as current list
			err := pkg.SetCurrentList(listName)
			if err != nil {
				fmt.Printf("Error setting current list: %v\n", err)
				return
			}
			
			// Create todo file if it doesn't exist
			if !pkg.TodoFileExists(listName) {
				err = pkg.CreateTodoFile(listName)
				if err != nil {
					fmt.Printf("Error creating todo file: %v\n", err)
					return
				}
				fmt.Printf("Created todo list '%s'\n", listName)
			} else {
				fmt.Printf("Switched to list '%s'\n", listName)
			}
			
			// Display current todos
			err = pkg.DisplayTodoList(listName)
			if err != nil {
				fmt.Printf("Error displaying todo list: %v\n", err)
				return
			}
		}
	},
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show history of completed todos across all lists",
	Long:  `Display a chronological history of all completed todos with timestamps, organized by date.`,
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		err := pkg.ShowHistory()
		if err != nil {
			fmt.Printf("Failed to show history: %v\n", err)
			return
		}
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Output comprehensive information about todo CLI for LLM assistants",
	Long:  `Outputs detailed information about the todo CLI structure, commands, and usage patterns designed for LLM assistants to understand how to use the tool effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`# Todo CLI - LLM Assistant Guide

## Overview
This is a CLI tool for managing todo lists. Each todo list is stored as a markdown file in the .todo directory.

## Core Concepts
- **Lists**: Each todo list is a separate markdown file
- **Storage**: Todo items stored in .todo/<list-name>.md files
- **Current List**: Track which list is currently active via .current-list file

## Available Commands

### 1. todo init
Initialize todo management in the current directory.
- Use when: Directory lacks .todo setup
- Creates: .todo directory for storing todo files

### 2. todo list [list-name]
Manage todo lists (create, switch, view, delete).
- 'todo list' - Show all lists with progress percentages
- 'todo list <name>' - Switch to or create list
- 'todo list --delete <name>' - Delete list (requires confirmation)

### 3. todo add <item>
Add todo item to current list.
- Takes: Single quoted string argument
- Example: todo add "Implement user authentication"

### 4. todo check <number>
Mark todo item as completed.
- Takes: Item number (1-based indexing)
- Example: todo check 1

### 5. todo uncheck <number>
Mark todo item as incomplete.
- Takes: Item number (1-based indexing)
- Example: todo uncheck 2

### 6. todo progress [list-name]
Show progress for lists.
- 'todo progress' - Current list progress
- 'todo progress <name>' - Specific list progress
- 'todo progress --all' - All lists progress

### 7. todo history
Show chronological history of completed todos across all lists.

### 8. todo edit
Open current list in your configured editor ($EDITOR).

### 9. todo version
Show CLI version.

## File Structure
` + "```" + `
project/
├── .todo/
│   ├── main.md
│   ├── feature-auth.md
│   └── bug-fixes.md
├── .current-list
└── [project files]
` + "```" + `

## Todo File Format
Standard markdown with checkboxes:
` + "```" + `
# Todo List for feature-auth

- [ ] Incomplete task
- [x] Completed task (completed: 2024-01-15 10:30)
` + "```" + `

## Common Workflows

### Starting New Work
1. 'todo list my-feature' (creates or switches to list)
2. 'todo add "First task"'
3. 'todo add "Second task"'
4. Work and mark complete: 'todo check 1'

### Checking Progress
- 'todo progress' (current list)
- 'todo list' (all lists overview)
- 'todo progress --all' (detailed all lists)

### Switching Between Lists
- 'todo list other-feature' (switches to different list)

### Cleanup
- 'todo list --delete completed-feature' (removes todo file)

## Error Handling
- Creates .todo directory automatically if missing
- Prevents deleting currently active list
- Validates item numbers for check/uncheck
- Confirms before deleting lists

## Tips for LLM Assistants
1. Always check current status with 'todo progress' or 'todo list'
2. Use descriptive names for lists (feature names, project areas)
3. Add todos in logical order (dependencies first)
4. Check progress regularly to track completion
5. Clean up completed lists to maintain organization

## Safety Features
- Delete confirmations
- Current list protection
- Automatic directory creation
- Timestamp tracking for completed items

This tool is designed for developers who want flexible todo management.
`)
	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the current todo list in your configured editor",
	Long:  `Open the current todo list file in your configured editor (set via $EDITOR environment variable).`,
	Run: func(cmd *cobra.Command, args []string) {
		if requiresInit() {
			return
		}
		
		currentList, err := pkg.GetCurrentList()
		if err != nil {
			fmt.Printf("Error getting current list: %v\n", err)
			return
		}
		
		err = pkg.EditTodoFile(currentList)
		if err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			return
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of todo CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("todo CLI v0.2.0")
	},
}

func init() {
	// Add the --all flag to progress command
	progressCmd.Flags().BoolP("all", "a", false, "Show progress for all features")
	
	// Add the --delete flag to list command
	listCmd.Flags().BoolP("delete", "d", false, "Delete the specified list")
	
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(uncheckCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	rootCmd.SetUsageTemplate(`Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}