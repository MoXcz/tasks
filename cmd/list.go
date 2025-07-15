/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/MoXcz/tasks/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
func newListCmd(storage *file.FileStorage) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("this command does not accept any arguments")
			}

			if err := (*storage).ListTasks(cmd.OutOrStdout()); err != nil {
				fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
				return fmt.Errorf("error listing tasks: %w", err)
			}
			return nil
		},
	}

	listCmd.Flags().BoolP("all", "a", false, "List all tasks including completed ones")
	if err := viper.BindPFlag("all", listCmd.Flags().Lookup("all")); err != nil {
		fmt.Fprintf(os.Stderr, "error binding --all flag: %v\n", err)
	}
	return listCmd
}
