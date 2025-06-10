/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all tasks",
	Long: `list all tasks
tasks list to list all tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("This command does not accept any arguments.")
			return
		}

		if err := storage.ListTasks(); err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
