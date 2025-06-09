/*
Copyright © 2025 Oscar Marquez
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/MoXcz/tasks/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var id int = 0

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

		if err := addTask(args[0]); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}
	},
}

func addTask(task string) error {
	taskFile := viper.GetString("file")
	file, err := helpers.LoadFile(taskFile)
	if err != nil {
		return fmt.Errorf("error loading task file: %w", err)
	}
	defer helpers.CloseFile(file)

	csvWriter := csv.NewWriter(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	// Write the header if the file is empty
	if fileInfo.Size() == 0 {
		if err := csvWriter.Write([]string{"ID", "Task", "CreatedAt", "IsComplete"}); err != nil {
			return fmt.Errorf("error writing header to file: %w", err)
		}
	}

	record := []string{strconv.Itoa(id), task, time.Now().Format("Mon Jan 2 15:04:05"), "false"}
	// TODO: use time.Parse(), read the time and then timediff.TimeDiff(time.Now().Add(time.Parse()))
	if err := csvWriter.Write(record); err != nil {
		return fmt.Errorf("error writing task to file: %w", err)
	}
	id++

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
