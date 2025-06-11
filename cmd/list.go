/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all tasks",
	Long: `list all tasks
tasks list to list all tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Fprintln(os.Stdin, "This command does not accept any arguments.")
			return
		}

		if err := storage.ListTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "List all tasks including completed ones")
	viper.BindPFlag("all", listCmd.Flags().Lookup("all"))
}
