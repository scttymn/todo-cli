package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/scttymn/todo-cli/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A CLI tool for managing branch-specific todo lists",
	Long:  `todo is a CLI tool that manages todo lists tied to git branches, helping you track tasks for each feature or project.`,
}

var startCmd = &cobra.Command{
	Use:   "start [feature-name]",
	Short: "Start working on a feature (create/switch to branch and manage todos)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]
		branchName := "feature/" + featureName
		
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
			fmt.Printf("Switched to existing branch '%s'\n", branchName)
		} else {
			// Branch doesn't exist, create it
			err = pkg.CreateBranch(branchName)
			if err != nil {
				fmt.Printf("Error creating branch %s: %v\n", branchName, err)
				return
			}
			fmt.Printf("Created and switched to branch '%s'\n", branchName)
		}
		
		// Ensure todo file exists
		if !pkg.TodoFileExists(featureName) {
			err = pkg.CreateTodoFile(featureName)
			if err != nil {
				fmt.Printf("Error creating todo file: %v\n", err)
				return
			}
			fmt.Printf("Initialized todo file: .todo/%s.md\n", featureName)
		}
		
		// Display current todos
		err = pkg.DisplayTodoList(featureName)
		if err != nil {
			fmt.Printf("Error displaying todo list: %v\n", err)
			return
		}
	},
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
	
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(uncheckCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}