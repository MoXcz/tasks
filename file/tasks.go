package file

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/viper"
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

	createdAt, err := time.Parse(time.RFC1123, record[2])
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

func printTasks(tasks []Task) {
	printAll := viper.GetBool("all")
	tabW := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	tabW.Write(fmt.Appendf(nil, "ID\t Task\t Created\t Done\n"))

	for _, task := range tasks {
		formattedCreatedAt := timediff.TimeDiff(task.CreatedAt)

		if !task.IsComplete {
			tabW.Write(fmt.Appendf(nil, "%d\t %s\t %s\t %t\n", task.ID, task.Task, formattedCreatedAt, task.IsComplete))
			continue
		}
		if printAll {
			tabW.Write(fmt.Appendf(nil, "%d\t %s\t %s\t %t\n", task.ID, task.Task, formattedCreatedAt, task.IsComplete))
		}
	}

	if err := tabW.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not print tasks: %v\n", err)
		return
	}
}
