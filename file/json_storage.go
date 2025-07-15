/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
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
	return nil
}

func (s *JSONStorage) DeleteTask(w io.Writer, id int) error {
	return nil
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
