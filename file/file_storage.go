/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"fmt"

	"github.com/spf13/viper"
)

type FileStorage interface {
	AddTask(task string) error
	ListTasks() error
	CompleteTask(id int) error
	DeleteTask(id int) error
}

func SelectStorage() (FileStorage, error) {
	storageType := viper.GetString("storage")
	path := viper.GetString("file")

	switch storageType {
	case "csv":
		storage, err := NewCSVStorage(path)
		if err != nil {
			return nil, fmt.Errorf("error creating CSV storage: %w", err)
		}
		return storage, nil
		// case "json"
		// case "sqlite":
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
