/*
Copyright © 2025 Oscar Marquez
*/
package file

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// Load file with exclusive lock.
// Check if the file exists, if not create it with read-write permissions
func LoadFile(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}

	return file, nil
}

// Release the lock on the file descriptor and close the file
func CloseFile(file *os.File) error {
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	return file.Close()
}

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
