package file

import (
	"fmt"
	"io"
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

func newTask(taskID, task, created, isComplete string) (Task, error) {
	// TODO: is it necessary to type check this?
	if taskID == "" {
		return Task{}, fmt.Errorf("taksID is empty: %s", taskID)
	}
	if task == "" {
		return Task{}, fmt.Errorf("task description is empty: %s", task)
	}
	if created == "" {
		return Task{}, fmt.Errorf("createdAt is empty: %s", created)
	}
	if isComplete == "" {
		return Task{}, fmt.Errorf("isComplete is empty: %s", created)
	}

	id, err := strconv.Atoi(taskID)
	if err != nil {
		return Task{}, fmt.Errorf("error converting ID to integer: %w", err)
	}

	createdAt, err := time.Parse(time.RFC1123, created)
	if err != nil {
		return Task{}, fmt.Errorf("error parsing created at time: %w", err)
	}

	return Task{
		ID:         id,
		Task:       task,
		CreatedAt:  createdAt,
		IsComplete: isComplete == "true", // any other value than "true" will default to "false"
	}, nil
}

func printTasks(w io.Writer, tasks []Task) {
	printAll := viper.GetBool("all")
	tabW := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)

	// TODO: perhaps it's better to return the error here?
	_, err := tabW.Write(fmt.Appendf(nil, "ID\t Task\t Created\t Done\n"))
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not write header to stdout: %v\n", err)
		return
	}

	for _, task := range tasks {
		formattedCreatedAt := timediff.TimeDiff(task.CreatedAt)

		if !task.IsComplete {
			_, err := tabW.Write(fmt.Appendf(nil, "%d\t %s\t %s\t %t\n", task.ID, task.Task, formattedCreatedAt, task.IsComplete))
			if err != nil {
				fmt.Fprintf(os.Stdout, "Could not write to stdout: %v\n", err)
				return
			}
			continue
		}
		if printAll {
			_, err := tabW.Write(fmt.Appendf(nil, "%d\t %s\t %s\t %t\n", task.ID, task.Task, formattedCreatedAt, task.IsComplete))
			if err != nil {
				fmt.Fprintf(os.Stdout, "Could not write to stdout: %v\n", err)
				return
			}
		}
	}

	if err := tabW.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not print tasks: %v\n", err)
		return
	}
}
