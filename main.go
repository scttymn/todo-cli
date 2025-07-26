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

func requiresGitSetup() bool {
	if !pkg.IsGitRepository() || !pkg.HasCommits() {
		fmt.Println("This directory is not set up for todo management.")
		fmt.Println("Run 'todo init' to initialize a todo-enabled git repository.")
		return true
	}
	return false
}

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI tool for managing branch-specific todo lists",
	Long:  `todo is a CLI tool that manages todo lists tied to git branches, helping you track tasks for each feature or project.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new todo-enabled git repository",
	Long:  `Initialize the current directory as a git repository with proper setup for todo management.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.InitTodoRepository()
		if err != nil {
			fmt.Printf("Failed to initialize todo repository: %v\n", err)
			return
		}
		
		fmt.Println("✅ Todo repository initialized successfully!")
		fmt.Println("You can now create todo lists with: todo list <name>")
	},
}


var addCmd = &cobra.Command{
	Use:   "add [todo-item]",
	Short: "Add a todo item to the current branch's list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresGitSetup() {
			return
		}
		
		todoItem := args[0]
		
		featureName, err := pkg.GetFeatureName()
		if err != nil {
			fmt.Printf("Error getting feature name: %v\n", err)
			return
		}
		
		err = pkg.AddTodoItem(featureName, todoItem)
		if err != nil {
			fmt.Printf("Error adding todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Added todo item to feature '%s': %s\n", featureName, todoItem)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [item-number]",
	Short: "Mark a todo item as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresGitSetup() {
			return
		}
		
		itemNumber := args[0]
		
		featureName, err := pkg.GetFeatureName()
		if err != nil {
			fmt.Printf("Error getting feature name: %v\n", err)
			return
		}
		
		itemID, err := strconv.Atoi(itemNumber)
		if err != nil {
			fmt.Printf("Invalid item number: %s\n", itemNumber)
			return
		}
		
		err = pkg.CheckTodoItem(featureName, itemID)
		if err != nil {
			fmt.Printf("Error checking todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Marked item %d as completed in feature '%s'\n", itemID, featureName)
	},
}

var uncheckCmd = &cobra.Command{
	Use:   "uncheck [item-number]",
	Short: "Mark a todo item as not completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresGitSetup() {
			return
		}
		
		itemNumber := args[0]
		
		featureName, err := pkg.GetFeatureName()
		if err != nil {
			fmt.Printf("Error getting feature name: %v\n", err)
			return
		}
		
		itemID, err := strconv.Atoi(itemNumber)
		if err != nil {
			fmt.Printf("Invalid item number: %s\n", itemNumber)
			return
		}
		
		err = pkg.UncheckTodoItem(featureName, itemID)
		if err != nil {
			fmt.Printf("Error unchecking todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Marked item %d as not completed in feature '%s'\n", itemID, featureName)
	},
}

var progressCmd = &cobra.Command{
	Use:   "progress [list-name]",
	Short: "Show progress for current list, specific list, or all lists",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresGitSetup() {
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
			featureName, err := pkg.GetFeatureName()
			if err != nil {
				fmt.Printf("Error getting feature name: %v\n", err)
				return
			}
			
			err = pkg.DisplayTodoList(featureName)
			if err != nil {
				fmt.Printf("Error displaying todo list: %v\n", err)
				return
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list [list-name]",
	Short: "Show all lists or switch to/create/delete a specific list",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if requiresGitSetup() {
			return
		}
		
		deleteFlag, _ := cmd.Flags().GetBool("delete")
		
		if deleteFlag {
			if len(args) == 0 {
				fmt.Println("Error: --delete requires a list name")
				return
			}
			
			listName := args[0]
			branchName := "feature/" + listName
			
			// Check if we're currently on the branch we're trying to delete
			currentBranch, err := pkg.GetCurrentBranch()
			if err != nil {
				fmt.Printf("Error getting current branch: %v\n", err)
				return
			}
			
			if currentBranch == branchName {
				fmt.Printf("Error: Cannot delete list '%s' because you are currently on it.\n", listName)
				fmt.Println("Switch to another list first (e.g., 'todo list main')")
				return
			}
			
			// Check if branch exists
			branchExists, err := pkg.BranchExists(branchName)
			if err != nil {
				fmt.Printf("Error checking if branch exists: %v\n", err)
				return
			}
			
			if !branchExists {
				fmt.Printf("List '%s' does not exist\n", listName)
				return
			}
			
			// Confirmation prompt
			fmt.Printf("Are you sure you want to delete list '%s'? This will remove both the branch and todo file. (y/N): ", listName)
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
			
			// Delete the branch
			err = pkg.DeleteBranch(branchName)
			if err != nil {
				fmt.Printf("Error deleting branch: %v\n", err)
				return
			}
			
			// Delete the todo file
			todoFile := pkg.GetTodoFilePath(listName)
			if err := os.Remove(todoFile); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: Could not delete todo file: %v\n", err)
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
			branchName := "feature/" + listName
			
			// Check for uncommitted changes before switching
			hasChanges, err := pkg.HasUncommittedChanges()
			if err != nil {
				fmt.Printf("Error checking for uncommitted changes: %v\n", err)
				return
			}
			
			if hasChanges {
				fmt.Println("⚠️  Warning: You have uncommitted changes that will be brought to the new branch.")
				fmt.Println("Consider committing or stashing your changes first.")
				fmt.Print("Do you want to continue? (y/N): ")
				
				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					return
				}
				
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					fmt.Println("Operation cancelled.")
					return
				}
			}
			
			// Check if branch exists
			branchExists, err := pkg.BranchExists(branchName)
			if err != nil {
				fmt.Printf("Error checking branch existence: %v\n", err)
				return
			}
			
			if branchExists {
				// Branch exists, just switch to it
				err = pkg.SwitchBranch(branchName)
				if err != nil {
					fmt.Printf("Error switching to branch %s: %v\n", branchName, err)
					return
				}
				fmt.Printf("Switched to existing list '%s'\n", listName)
			} else {
				// Branch doesn't exist, create it
				err = pkg.CreateBranch(branchName)
				if err != nil {
					fmt.Printf("Error creating branch %s: %v\n", branchName, err)
					return
				}
				fmt.Printf("Created and switched to list '%s'\n", listName)
			}
			
			// Ensure todo file exists
			if !pkg.TodoFileExists(listName) {
				err = pkg.CreateTodoFile(listName)
				if err != nil {
					fmt.Printf("Error creating todo file: %v\n", err)
					return
				}
				fmt.Printf("Initialized todo file: .todo/%s.md\n", listName)
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
		if requiresGitSetup() {
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
This is a Git-integrated CLI tool for managing branch-specific todo lists. Each todo list is tied to a Git branch and stored as a markdown file.

## Core Concepts
- **Lists**: Each todo list corresponds to a Git branch (pattern: feature/<list-name>)
- **Storage**: Todo items stored in .todo/<list-name>.md files
- **Branch Integration**: Switching lists switches Git branches
- **Safety**: Warns about uncommitted changes before switching

## Available Commands

### 1. todo init
Initialize a new todo-enabled git repository.
- Use when: Directory is not a git repo or lacks proper setup
- Creates: Git repository with todo infrastructure

### 2. todo list [list-name]
Manage todo lists (create, switch, view, delete).
- 'todo list' - Show all lists with progress percentages
- 'todo list <name>' - Switch to or create list (creates feature/<name> branch)
- 'todo list --delete <name>' - Delete list and branch (requires confirmation)

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

### 8. todo version
Show CLI version.

## File Structure
` + "```" + `
project/
├── .todo/
│   ├── feature-name.md
│   └── another-list.md
└── [project files]
` + "```" + `

## Todo File Format
Standard markdown with checkboxes:
` + "```" + `
# Todo List for feature-name

- [ ] Incomplete task
- [x] Completed task
` + "```" + `

## Common Workflows

### Starting New Feature Work
1. 'todo list feature-name' (creates branch + todo file)
2. 'todo add "First task"'
3. 'todo add "Second task"'
4. Work and mark complete: 'todo check 1'

### Checking Progress
- 'todo progress' (current list)
- 'todo list' (all lists overview)
- 'todo progress --all' (detailed all lists)

### Switching Between Features
- 'todo list other-feature' (switches branch + shows todos)
- Warns if uncommitted changes exist

### Cleanup
- 'todo list --delete completed-feature' (removes branch + file)

## Error Handling
- Requires git repository (suggests 'todo init' if missing)
- Warns about uncommitted changes before branch switches
- Prevents deleting currently active list
- Validates item numbers for check/uncheck

## Tips for LLM Assistants
1. Always check current status with 'todo progress' or 'todo list'
2. Use descriptive names for lists (feature names, not generic terms)
3. Add todos in logical order (dependencies first)
4. Check progress regularly to track completion
5. Clean up completed lists to maintain organization
6. Remember that switching lists switches Git branches

## Safety Features
- Git repository validation
- Uncommitted changes warnings
- Delete confirmations
- Branch existence checks
- Current branch protection

This tool is designed for developers who want to track feature-specific tasks while maintaining clean Git branch organization.
`)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of todo CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("todo CLI v0.1.0")
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
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}