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

func (s *CSVStorage) ListTasks() ([]string, error) {
	return nil, nil
}

func (s *CSVStorage) CompleteTask(id int) error {
	return nil
}
