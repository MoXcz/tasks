/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type CSVStorage struct {
	filepath string
}

func NewCSVStorage(filepath string) *CSVStorage {
	return &CSVStorage{
		filepath: filepath,
	}
}

func (s *CSVStorage) AddTask(task string) error {
	file, err := LoadFile(s.filepath)
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

	record := []string{strconv.Itoa(lastID + 1), task, time.Now().UTC().Format(time.RFC1123), "false"}
	if err := csvWriter.Write(record); err != nil {
		return fmt.Errorf("error writing task to file: %w", err)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func (s *CSVStorage) ListTasks(w io.Writer) error {
	tasks, err := readTasksCSV(s.filepath)
	if err != nil {
		return err // Error already formatted in readTasksCSV
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found in the file: header only")
		return nil
	}

	fmt.Fprintln(w, "Total tasks:", len(tasks))
	printTasks(w, tasks)

	return nil
}

func (s *CSVStorage) CompleteTask(w io.Writer, id int) error {
	if id <= 0 {
		return fmt.Errorf("task ID must be greater than 0, got %d", id)
	}

	found := false
	tasks, err := readTasksCSV(s.filepath)
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
			fmt.Fprintln(w, "Completing task:", task.Task)
			tasks[i].IsComplete = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return writeTasksCSV(s.filepath, tasks)
}

func (s *CSVStorage) DeleteTask(w io.Writer, id int) error {
	force := viper.GetBool("force")
	if id <= 0 {
		return fmt.Errorf("task ID must be greater than 0 %d", id)
	}

	tasks, err := readTasksCSV(s.filepath)
	if err != nil {
		return err // Error already formatted in readTasksCSV
	}

	var found bool

	for i, task := range tasks {
		if task.ID == id && !task.IsComplete {
			for {
				if !force {
					fmt.Printf("Are you sure you want to delete this uncompleted task ([y]es | [n]o)? ")
					reader := bufio.NewReader(os.Stdin)
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("An error occured while reading input. Please try again", err)
						continue
					}

					input = strings.TrimSuffix(input, "\n") // remove trailing \n

					if input == "no" || input == "n" {
						return nil
					}

					if input == "yes" || input == "y" {
						break // enter next conditional and delete the task
					}
				} else {
					break
				}
			}
		}

		if task.ID == id {
			fmt.Fprintln(w, "Deleting task:", task.Task)
			tasks = slices.Delete(tasks, i, i+1) // delete current task
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return writeTasksCSV(s.filepath, tasks)
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

	// TODO: check whether if longer returning an error broke something
	if len(records) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	for _, record := range records {
		if record[0] == "ID" {
			continue // Skip the header row
		}
		task, err := newTask(record[0], record[1], record[2], record[3])
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func writeTasksCSV(path string, tasks []Task) error {
	file, err := LoadFile(path)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	os.Truncate(path, 0)
	csvWriter := csv.NewWriter(file)
	if err := csvWriter.Write([]string{"ID", "Task", "CreatedAt", "IsComplete"}); err != nil {
		return fmt.Errorf("error writing header to file: %w", err)
	}

	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Task,
			task.CreatedAt.Format(time.RFC1123),
			strconv.FormatBool(task.IsComplete),
		}
		// TODO: check this, because if it fails it will remove the file contents, or fill it halfway at the point it errors out
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("error writing task to file: %w", err)
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
