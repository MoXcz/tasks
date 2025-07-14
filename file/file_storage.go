/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"fmt"
	"io"
)

type FileStorage interface {
	AddTask(task string) error
	ListTasks(w io.Writer) error
	CompleteTask(w io.Writer, id int) error
	DeleteTask(w io.Writer, id int) error
}

func SelectStorage(path, storageType string) (FileStorage, error) {
	switch storageType {
	case "csv":
		return NewCSVStorage(path + "." + storageType), nil
	case "json":
		return NewJSONStorage(path + "." + storageType), nil
		// case "sqlite":
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
