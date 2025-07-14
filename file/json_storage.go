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
	var tasks []Task
	file, err := LoadFile(s.filepath)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}
	defer CloseFile(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Println(string(fileContent))

	lastID := 0
	if fileInfo.Size() != 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return fmt.Errorf("error unmarshalling file: %w", err)
		}
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
	return nil
}

func (s *JSONStorage) CompleteTask(w io.Writer, id int) error {
	return nil
}

func (s *JSONStorage) DeleteTask(w io.Writer, id int) error {
	return nil
}
