package file

import (
	"fmt"
	"strconv"
	"time"
)

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
