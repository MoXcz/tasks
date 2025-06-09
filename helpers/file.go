package helpers

import (
	"fmt"
	"os"
	"syscall"
)

func LoadFile(filepath string) (*os.File, error) {
	// Check if the file exists, if not create it with read-write permissions
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

func CloseFile(file *os.File) error {
	// Release the lock on the file descriptor
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	return file.Close()
}
