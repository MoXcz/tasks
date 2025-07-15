/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MoXcz/tasks/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
func newDeleteCmd(storage *file.FileStorage) *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete a task",
		Long: `delete a task
tasks delete <task ID> to delete a task from the list`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Please provide exactly one task ID to delete.")
				return
			}

			ID := args[0]
			taskID, err := strconv.Atoi(ID)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Invalid task ID: %s.\n", ID)
				return
			}

			if err := (*storage).DeleteTask(cmd.OutOrStdout(), taskID); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error deleting task: %v\n", err)
				return
			}
		},
	}

	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
	if err := viper.BindPFlag("force", deleteCmd.Flags().Lookup("force")); err != nil {
		fmt.Fprintf(os.Stderr, "error binding --force flag: %v\n", err)
	}
	return deleteCmd
}
