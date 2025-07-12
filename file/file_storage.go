/*
Copyright © 2025 Oscar Marquez
*/
package file

import "fmt"

type FileStorage interface {
	AddTask(task string) error
	ListTasks() error
	CompleteTask(id int) error
	DeleteTask(id int) error
}

func SelectStorage(path, storageType string) (FileStorage, error) {
	switch storageType {
	case "csv":
		return NewCSVStorage(path), nil
	case "json":
		return NewJSONStorage(path), nil
		// case "sqlite":
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
