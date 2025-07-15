/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type JSONStorage struct {
	filepath string
}

func NewJSONStorage(filepath string) *JSONStorage {
	return &JSONStorage{
		filepath: filepath,
	}
}

func (s *JSONStorage) AddTask(task string) error {
	tasks, err := readTasksJSON(s.filepath)
	if err != nil {
		return err
	}

	lastID := 0
	if len(tasks) != 0 {
		lastID = tasks[len(tasks)-1].ID // should be last task
	}

	tasks = append(tasks, Task{
		ID:         lastID,
		Task:       task,
		CreatedAt:  time.Now().UTC(),
		IsComplete: false,
	})

	data, err := json.Marshal(&tasks)
	if err != nil {
		return fmt.Errorf("could not marshal tasks: %w", err)
	}

	os.WriteFile(s.filepath, data, 0644)
	return nil
}

func (s *JSONStorage) ListTasks(w io.Writer) error {
	tasks, err := readTasksJSON(s.filepath)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found in the file: header only")
		return nil
	}

	fmt.Fprintln(w, "Total tasks:", len(tasks))
	printTasks(w, tasks)

	return nil
}

func (s *JSONStorage) CompleteTask(w io.Writer, id int) error {
	if id <= 0 {
		return fmt.Errorf("task ID must be greater than 0, got %d", id)
	}

	found := false
	tasks, err := readTasksJSON(s.filepath)
	if err != nil {
		return err
	}

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
	return writeTasksJSON(s.filepath, tasks)
}

func (s *JSONStorage) DeleteTask(w io.Writer, id int) error {
	force := viper.GetBool("force")
	if id <= 0 {
		return fmt.Errorf("task ID must be greater than 0 %d", id)
	}

	tasks, err := readTasksJSON(s.filepath)
	if err != nil {
		return err
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

	err = writeTasksJSON(s.filepath, tasks)
	return err
}

func readTasksJSON(path string) ([]Task, error) {
	var tasks []Task
	file, err := LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if fileInfo.Size() == 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return nil, fmt.Errorf("error unmarshalling file: %w", err)
		}
		return []Task{}, nil // no tasks, does nil work here?
	}
	return tasks, nil
}

func writeTasksJSON(path string, tasks []Task) error {
	file, err := LoadFile(path)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	os.Truncate(path, 0)

	// TODO: check this, because if it fails it will remove the file contents, or fill it halfway at the point it errors out
	data, err := json.Marshal(&tasks)
	if err != nil {
		return fmt.Errorf("could not marshal tasks: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
