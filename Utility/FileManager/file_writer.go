package FileManager

import (
	"fmt"
	"io/fs"
	"os"
)

// FileWriter represents a file writer that can be used to write to files.
type FileWriter struct {
	// loc is the location of the file.
	loc string

	// file is the file that is being written to.
	file *os.File
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (fw *FileWriter) GetLocation() string {
	return fw.loc
}

// NewFileWriter creates a new FileWriter with the given location.
//
// Parameters:
//   - loc: A string representing the location of the file.
//
// Returns:
//   - *FileWriter: A pointer to the newly created FileWriter.
func NewFileWriter(loc string) *FileWriter {
	return &FileWriter{
		loc:  loc,
		file: nil,
	}
}

// Exists checks if the file exists at the location.
//
// Parameters:
//   - loc: A string representing the location of the file.
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (fw *FileWriter) Exists() (bool, error) {
	if fw.file != nil {
		return true, nil
	}

	return fileExists(fw.loc)
}

// Create creates a new file at the location; including any necessary directories.
//
// Parameters:
//   - perm: An fs.FileMode representing the permissions to set on the file.
//
// Returns:
//   - error: An error if one occurred while creating the file.
//
// Behaviors:
//   - If the file already exists, the function does nothing.
//   - Once the file is opened, it is kept open until the FileWriter is closed.
func (fw *FileWriter) Create(perm fs.FileMode) error {
	if fw.file != nil {
		return nil
	}

	file, err := CreateAll(fw.loc, perm)
	if err != nil {
		return err
	}

	fw.file = file

	return nil
}

// Open opens the file for writing.
//
// Returns:
//   - *os.File: A pointer to the opened file.
//   - error: An error if one occurred while opening the file.
//
// Behaviors:
//   - The file is opened in append mode and write-only.
//   - If the file is already open, the function does nothing.
//   - If the file does not exist, the function creates the file.
func (fw *FileWriter) Open() error {
	if fw.file != nil {
		return nil
	}

	file, err := os.OpenFile(fw.loc, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fw.file = file

	return nil
}

// Close closes the file that is being written to.
//
// Returns:
//   - error: An error if one occurred while closing the file.
func (fw *FileWriter) Close() error {
	if fw.file == nil {
		return nil
	}

	err := fw.file.Close()
	if err != nil {
		return err
	}

	fw.file = nil

	return nil
}

// AppendLine appends a line of content to the file.
//
// Parameters:
//   - content: A string representing the content to append to the file.
//
// Returns:
//   - error: An error if one occurred while writing to the file.
//
// Errors:
//   - *ErrFileNotOpen: If the file is not open.
//   - error: If any other case.
func (fw *FileWriter) AppendLine(content string) error {
	if fw.file == nil {
		return NewErrFileNotOpen()
	}

	_, err := fw.file.WriteString(content + "\n")
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

// EmptyLine appends an empty line to the file.
//
// Returns:
//   - error: An error if one occurred while writing to the file.
//
// Errors:
//   - *ErrFileNotOpen: If the file is not open.
//   - error: If any other case.
func (fw *FileWriter) EmptyLine() error {
	if fw.file == nil {
		return NewErrFileNotOpen()
	}

	_, err := fw.file.WriteString("\n")
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

// Clear clears the contents of the file.
//
// Returns:
//   - error: An error if one occurred while truncating the file.
//
// Errors:
//   - *ErrFileNotOpen: If the file is not open.
//   - *os.PathError: If any other case.
func (fw *FileWriter) Clear() error {
	if fw.file == nil {
		return NewErrFileNotOpen()
	}

	err := fw.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("could not clear file: %w", err)
	}

	return nil
}

// Write implements the io.Writer interface for the FileWriter.
func (fw *FileWriter) Write(p []byte) (n int, err error) {
	if fw.file == nil {
		return 0, NewErrFileNotOpen()
	}

	return fw.file.Write(p)
}
