/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"

	"github.com/MoXcz/tasks/file"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
func newAddCmd(storage *file.FileStorage) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add a new task",
		Long: `add a new task
tasks add <task description> to add a new task`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Fprintln(cmd.OutOrStdout(), "Please provide a task description.")
				return
			}

			// TODO: Manage task description with spaces
			if len(args) > 1 {
				fmt.Fprintf(cmd.OutOrStdout(), "Please provide the task description surrounded by quotes\n\nExample: task add \"Go to the gym\"\n")
				return
			}

			task := args[0]
			if err := (*storage).AddTask(task); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error adding task: %v\n", err)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Task added successfully: %s\n", task)
		},
	}
}
