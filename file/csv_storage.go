/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"encoding/csv"
	"fmt"
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

	record := []string{strconv.Itoa(lastID + 1), task, time.Now().Format("Mon Jan 2 15:04:05"), "false"}
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
	file, err := LoadFile(s.path)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("No tasks found in the file")
		return nil
	}

	var tasks []Task
	for _, record := range records {
		if record[0] == "ID" {
			continue // Skip the header row
		}
		task, err := newTask(record)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found in the file: header only")
		return nil
	}

	fmt.Println("Total tasks:", len(tasks))

	for _, task := range tasks {
		fmt.Printf("ID: %d, Task: %s, CreatedAt: %s, IsComplete: %t\n", task.ID, task.Task, task.CreatedAt.Format("Mon Jan 2 15:04:05"), task.IsComplete)
	}
	return nil
}

func (s *CSVStorage) CompleteTask(id int) error {
	return nil
}

type Task struct {
	ID         int
	Task       string
	CreatedAt  time.Time
	IsComplete bool
}

func newTask(record []string) (Task, error) {
	if len(record) < 4 {
		return Task{}, fmt.Errorf("record does not contain enough fields: %v", record)
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return Task{}, fmt.Errorf("error converting ID to integer: %w", err)
	}

	task := record[1]

	createdAt, err := time.Parse("Mon Jan 2 15:04:05", record[2])
	if err != nil {
		return Task{}, fmt.Errorf("error parsing created at time: %w", err)
	}

	isComplete := record[3]

	return Task{
		ID:         id,
		Task:       task,
		CreatedAt:  createdAt,
		IsComplete: isComplete == "true",
	}, nil
}
