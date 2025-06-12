/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type CSVStorage struct {
	path string
}

func NewCSVStorage(path string) (*CSVStorage, error) {
	taskFile := viper.GetString("file")
	return &CSVStorage{
		path: taskFile,
	}, nil
}

func (s *CSVStorage) AddTask(task string) error {
	file, err := LoadFile(s.path)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	csvWriter := csv.NewWriter(file)
	// Write the header if the file is empty
	if fileInfo.Size() == 0 {
		if err := csvWriter.Write([]string{"ID", "Task", "CreatedAt", "IsComplete"}); err != nil {
			return fmt.Errorf("error writing header to file: %w", err)
		}
	}

	csvReader := csv.NewReader(file)
	// Read the existing records to find the next ID
	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading existing records: %w", err)
	}

	lastID := 0 // if len(records) < 0, lastID will be 1 (+1 below)
	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		if lastRecord[0] != "ID" { // Skip header
			lastID, err = strconv.Atoi(lastRecord[0])
			if err != nil {
				return fmt.Errorf("error converting last ID to integer: %w", err)
			}
		}
	}

	record := []string{strconv.Itoa(lastID + 1), task, time.Now().Format(time.RFC1123), "false"}
	// TODO: use time.Parse(), read the time and then timediff.TimeDiff(time.Now().Add(time.Parse()))
	if err := csvWriter.Write(record); err != nil {
		return fmt.Errorf("error writing task to file: %w", err)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func (s *CSVStorage) ListTasks() error {
	tasks, err := readTasksCSV(s.path)
	if err != nil {
		return err // Error already formatted in readTasksCSV
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found in the file: header only")
		return nil
	}

	fmt.Println("Total tasks:", len(tasks))
	printTasks(tasks)

	return nil
}

func (s *CSVStorage) CompleteTask(id int) error {
	if id <= 0 {
		return fmt.Errorf("task ID must be greater than 0, got %d", id)
	}

	found := false
	tasks, err := readTasksCSV(s.path)
	if err != nil {
		return err
	}

	// I'd think it's necessary to iterate through the list of tasks to make sure
	// both that the task exists and that it's not already completed
	for i, task := range tasks {
		if task.ID == id && task.IsComplete {
			return fmt.Errorf("task with ID %d is already completed", id)
		}
		if task.ID == id {
			fmt.Println("Completing task:", task.Task)
			tasks[i].IsComplete = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	file, err := LoadFile(s.path)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	os.Truncate(s.path, 0)
	csvWriter := csv.NewWriter(file)
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Task,
			task.CreatedAt.Format(time.RFC1123),
			strconv.FormatBool(task.IsComplete),
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("error writing task to file: %w", err)
		}
	}

	csvWriter.Flush()
	return nil
}

func readTasksCSV(path string) ([]Task, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("No tasks found in the file")
		return nil, nil
	}
	var tasks []Task
	for _, record := range records {
		if record[0] == "ID" {
			continue // Skip the header row
		}
		task, err := newTask(record)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
