/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task",
	Long: `add a new task
tasks add <task description> to add a new task`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Manage task description with spaces
		if len(args) < 1 {
			fmt.Println("Please provide a task description.")
			return
		}

		task := args[0]
		if err := storage.AddTask(task); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}
		fmt.Printf("Task added successfully: %s\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
