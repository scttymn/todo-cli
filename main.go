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

var initCmd = &cobra.Command{
	Use:   "init [branch-name]",
	Short: "Create a new branch and initialize its todo list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		
		err := pkg.CreateBranch(branchName)
		if err != nil {
			fmt.Printf("Error creating branch %s: %v\n", branchName, err)
			return
		}
		
		err = pkg.SwitchBranch(branchName)
		if err != nil {
			fmt.Printf("Error switching to branch %s: %v\n", branchName, err)
			return
		}
		
		fmt.Printf("Created and switched to branch '%s'\n", branchName)
		
		err = pkg.CreateTodoFile(branchName)
		if err != nil {
			fmt.Printf("Error creating todo file: %v\n", err)
			return
		}
		
		fmt.Printf("Initialized todo file: .todo/%s.md\n", branchName)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [todo-item]",
	Short: "Add a todo item to the current branch's list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todoItem := args[0]
		
		branch, err := pkg.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			return
		}
		
		err = pkg.AddTodoItem(branch, todoItem)
		if err != nil {
			fmt.Printf("Error adding todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Added todo item to branch '%s': %s\n", branch, todoItem)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [item-number]",
	Short: "Mark a todo item as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		itemNumber := args[0]
		
		branch, err := pkg.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			return
		}
		
		itemID, err := strconv.Atoi(itemNumber)
		if err != nil {
			fmt.Printf("Invalid item number: %s\n", itemNumber)
			return
		}
		
		err = pkg.CheckTodoItem(branch, itemID)
		if err != nil {
			fmt.Printf("Error checking todo item: %v\n", err)
			return
		}
		
		fmt.Printf("Marked item %d as completed in branch '%s'\n", itemID, branch)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current todo progress for the active branch",
	Run: func(cmd *cobra.Command, args []string) {
		branch, err := pkg.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			return
		}
		
		err = pkg.DisplayTodoList(branch)
		if err != nil {
			fmt.Printf("Error displaying todo list: %v\n", err)
			return
		}
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch [branch-name]",
	Short: "Switch to a branch and show its todo list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		err := pkg.SwitchBranch(branchName)
		if err != nil {
			fmt.Printf("Error switching to branch %s: %v\n", branchName, err)
			return
		}
		fmt.Printf("Switched to branch '%s'\n", branchName)
		
		err = pkg.DisplayTodoList(branchName)
		if err != nil {
			fmt.Printf("Error displaying todo list: %v\n", err)
			return
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
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}