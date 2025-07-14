/*
Copyright © 2025 Oscar Marquez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/MoXcz/tasks/file"
	"github.com/MoXcz/tasks/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg config.Config) *cobra.Command {
	var storage file.FileStorage
	rootCmd := &cobra.Command{
		Use:   "tasks",
		Short: "A task management CLI that tries to mimic the common TODO web application",
		Long: `tasks is a CLI tool to manage your tasks, that's it.
tasks add <task> to add a new task
tasks complete <task id> to mark a task as completed
tasks list to list all tasks
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			storage, err = file.SelectStorage(cfg.Filepath, cfg.Storage)
			if err != nil {
				// TODO: most probably an invalid storage type, verify it
				// This could also happen when the config file and the tasks are on the same dir
				return fmt.Errorf("error selecting storage: %w", err)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&cfg.Filepath, "config", cfg.Filepath, "config file")
	rootCmd.PersistentFlags().StringVar(&cfg.Storage, "storage", cfg.Storage, "storage type")

	rootCmd.AddCommand(newListCmd(&storage))
	rootCmd.AddCommand(newDeleteCmd(&storage))
	rootCmd.AddCommand(newCompleteCmd(&storage))
	rootCmd.AddCommand(newAddCmd(&storage))

	return rootCmd
}
