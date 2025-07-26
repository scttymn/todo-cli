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

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI tool for managing branch-specific todo lists",
	Long:  `todo is a CLI tool that manages todo lists tied to git branches, helping you track tasks for each feature or project.`,
}


var addCmd = &cobra.Command{
	Use:   "add [todo-item]",
	Short: "Add a todo item to the current branch's list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
	Use:   "progress",
	Short: "Show progress for current feature or all features",
	Run: func(cmd *cobra.Command, args []string) {
		showAll, _ := cmd.Flags().GetBool("all")
		
		if showAll {
			err := pkg.ListAllFeatures()
			if err != nil {
				fmt.Printf("Error showing progress: %v\n", err)
				return
			}
		} else {
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
	
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(uncheckCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}