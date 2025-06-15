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
		storage := NewCSVStorage(path)
		return storage, nil
		// case "json"
		// case "sqlite":
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
