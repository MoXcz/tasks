/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "complete a task",
	Long: `complete a task
tasks complete <task ID> to mark a task as completed`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stdin, "Please provide exactly one task ID to complete.")
			return
		}

		ID := args[0]
		taskID, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid task ID: %s. Please provide a valid integer ID.\n", ID)
			return
		}

		if err := storage.CompleteTask(taskID); err != nil {
			fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
