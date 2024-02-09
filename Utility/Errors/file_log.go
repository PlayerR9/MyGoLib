package Errors

import (
	"fmt"
	"log"
	"os"
)

// FileLogger is a struct that represents a file logger.
// It must be closed after use to release the file resources.
type FileLogger struct {
	// The name of the file to log to.
	fileName string

	// The file to log to.
	file *os.File

	// The logger to write to the file.
	*log.Logger
}

// Close is a method of FileLogger that closes the file and releases the resources.
func (fl *FileLogger) Close() {
	if fl.file != nil {
		fl.file.Close()
		fl.file = nil
	}

	fl.Logger = nil
}

// NewFileLogger is a function that creates a new FileLogger instance.
// It creates a new file with the given file path and returns a pointer to the
// FileLogger instance.
//
// Parameters:
//
//   - filePath: The path of the file to log to.
//
// Returns:
//
//   - *FileLogger: A pointer to the new FileLogger instance.
//   - error: An error if the file cannot be created or opened for writing.
//
// The file name is created by appending the ".log.md" extension to the file path.
func NewFileLogger(filePath string) (*FileLogger, error) {
	fl := &FileLogger{
		fileName: fmt.Sprintf("%s.log.md", filePath),
	}

	// Open the file for writing
	var err error

	fl.file, err = os.OpenFile(fl.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot create or open file: %v", err)
	}

	fl.Logger = log.New(fl.file, "", log.LstdFlags|log.Llongfile)

	return fl, nil
}

// GetFileName is a method of FileLogger that returns the name of the file to log to.
//
// Returns:
//
//   - string: The name of the file to log to.
func (fl *FileLogger) GetFileName() string {
	return fl.fileName
}
