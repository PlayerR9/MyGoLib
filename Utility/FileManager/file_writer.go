package FileManager

import (
	"fmt"
	"os"
)

// FileWriter represents a file writer that can be used to write to files.
type FileWriter struct {
	*fileHandler
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
	return &FileWriter{
		fileHandler: newFileHandler(loc, dirPerm, filePerm, os.O_APPEND|os.O_WRONLY),
	}
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
		err := fw.fileHandler.Create()
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
		err := fw.fileHandler.Create()
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
		err := fw.fileHandler.Create()
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
		err := fw.fileHandler.Create()
		if err != nil {
			return 0, fmt.Errorf("could not create file: %w", err)
		}
	}

	return fw.file.Write(p)
}
