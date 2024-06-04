package FileManager

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter represents a file writer that can be used to write to files.
type FileWriter struct {
	// loc is the location of the file.
	loc string

	// file is the file that is being read from or written to.
	file *os.File

	// dirPerm is the permissions of the directory.
	dirPerm os.FileMode

	// filePerm is the permissions of the file.
	filePerm os.FileMode

	// flag is the flag to use when opening the file.
	flag int
}

// NewFileWriter creates a new FileWriter with the given location.
//
// Parameters:
//   - loc: A string representing the location of the file.
//   - dirPerm: An os.FileMode representing the permissions to set on the directories.
//   - filePerm: An os.FileMode representing the permissions to set on the file.
//
// Returns:
//   - *FileWriter: A pointer to the newly created FileWriter.
func NewFileWriter(loc string, dirPerm, filePerm os.FileMode) *FileWriter {
	if dirPerm == 0 {
		dirPerm = os.ModePerm
	}

	if filePerm == 0 {
		filePerm = os.ModePerm
	}

	return &FileWriter{
		loc:      loc,
		dirPerm:  dirPerm,
		filePerm: filePerm,
		flag:     os.O_APPEND | os.O_WRONLY,
	}
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (fw *FileWriter) GetLocation() string {
	return fw.loc
}

// GetPermissions returns the permissions of the file.
//
// Returns:
//   - [2]os.FileMode: An array of os.FileMode representing the permissions of the directory and file.
func (fw *FileWriter) GetPermissions() [2]os.FileMode {
	return [2]os.FileMode{fw.dirPerm, fw.filePerm}
}

// Close closes the file if it is open.
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

// Create creates the file at the location.
//
// Returns:
//   - error: An error if one occurred while creating the file.
//
// Behaviors:
//   - If the file already exists, it closes the previous file and creates a new one.
//   - Once the file is opened, it is kept open until the FileManager is closed.
func (fw *FileWriter) Create() error {
	if fw.file != nil {
		err := fw.file.Close()
		if err != nil {
			return err
		}
	}

	dir := filepath.Dir(fw.loc)

	err := os.MkdirAll(dir, fw.dirPerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fw.loc, os.O_CREATE|fw.flag, fw.filePerm)
	if err != nil {
		return err
	}

	fw.file = file

	return nil
}

// Exists checks if a file exists at the location of the FileWriter.
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (fw *FileWriter) Exists() (bool, error) {
	_, err := os.Stat(fw.loc)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Open opens the file at the location.
//
// Returns:
//   - error: An error if one occurred while opening the file.
//
// Behaviors:
//   - If the file is already open, the function closes the file and opens it again.
func (fw *FileWriter) Open() error {
	if fw.file != nil {
		err := fw.file.Close()
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(fw.loc, fw.flag, fw.filePerm)
	if err != nil {
		return err
	}

	fw.file = file

	return nil
}

// Delete deletes the file at the location.
//
// Returns:
//   - error: An error if one occurred while deleting the file.
func (fw *FileWriter) Delete() error {
	if fw.file != nil {
		err := fw.file.Close()
		if err != nil {
			return err
		}
	}

	return os.Remove(fw.loc)
}

// ChangePath changes the path of the file.
//
// Parameters:
//   - path: A string representing the new path of the file.
func (fw *FileWriter) ChangePath(path string) {
	fw.loc = path
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
		// Create the file if it does not exist
		err := fw.Create()
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		return nil
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
		// Create the file if it does not exist
		err := fw.Create()
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		return nil
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
		// Create the file if it does not exist
		err := fw.Create()
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		return nil
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
		// Create the file if it does not exist
		err := fw.Create()
		if err != nil {
			return 0, fmt.Errorf("could not create file: %w", err)
		}
	}

	return fw.file.Write(p)
}
