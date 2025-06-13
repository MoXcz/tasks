/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task",
	Long: `add a new task
tasks add <task description> to add a new task`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stdout, "Please provide a task description.")
			return
		}

		// TODO: Manage task description with spaces
		if len(args) > 1 {
			fmt.Fprintf(os.Stdout, "Please provide the task description surrounded by quotes\n\nExample: task add \"Go to the gym\"\n")
			return
		}

		task := args[0]
		if err := storage.AddTask(task); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding task: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Task added successfully: %s\n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
