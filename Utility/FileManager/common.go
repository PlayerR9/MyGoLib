package FileManager

import (
	"os"
	"path/filepath"
)

// fileHandler represents a file handler that can be used to read from and write to files.
type fileHandler struct {
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

// newFileHandler creates a new fileHandler.
//
// Parameters:
//   - loc: A string representing the location of the file.
//   - dirPerm: An os.FileMode representing the permissions to set on the directories.
//   - filePerm: An os.FileMode representing the permissions to set on the file.
//   - flag: An int representing the flag to use when opening the file.
//
// Returns:
//   - *fileHandler: A pointer to the newly created fileHandler.
func newFileHandler(loc string, dirPerm, filePerm os.FileMode, flag int) *fileHandler {
	return &fileHandler{
		loc:      loc,
		file:     nil,
		dirPerm:  dirPerm,
		filePerm: filePerm,
		flag:     flag,
	}
}

// GetLocation returns the location of the file.
//
// Returns:
//   - string: The location of the file.
func (fh *fileHandler) GetLocation() string {
	return fh.loc
}

// GetPermissions returns the permissions of the file.
//
// Returns:
//   - [2]os.FileMode: An array of os.FileMode representing the permissions of the directory and file.
func (fh *fileHandler) GetPermissions() [2]os.FileMode {
	return [2]os.FileMode{fh.dirPerm, fh.filePerm}
}

// Close closes the file if it is open.
//
// Returns:
//   - error: An error if one occurred while closing the file.
func (fh *fileHandler) Close() error {
	if fh.file == nil {
		return nil
	}

	err := fh.file.Close()
	if err != nil {
		return err
	}

	fh.file = nil

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
func (fh *fileHandler) Create() error {
	if fh.file != nil {
		err := fh.file.Close()
		if err != nil {
			return err
		}
	}

	dir := filepath.Dir(fh.loc)

	err := os.MkdirAll(dir, fh.dirPerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fh.loc, os.O_CREATE|fh.flag, fh.filePerm)
	if err != nil {
		return err
	}

	fh.file = file

	return nil
}

// Exists checks if a file exists at the location of the fileHandler.
//
// Returns:
//   - bool: A boolean indicating whether the file exists.
//   - error: An error if one occurred while checking the file.
func (fh *fileHandler) Exists() (bool, error) {
	_, err := os.Stat(fh.loc)
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
func (fh *fileHandler) Open() error {
	if fh.file != nil {
		err := fh.file.Close()
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(fh.loc, fh.flag, fh.filePerm)
	if err != nil {
		return err
	}

	fh.file = file

	return nil
}

// Delete deletes the file at the location.
//
// Returns:
//   - error: An error if one occurred while deleting the file.
func (fh *fileHandler) Delete() error {
	if fh.file != nil {
		err := fh.file.Close()
		if err != nil {
			return err
		}
	}

	return os.Remove(fh.loc)
}

// ChangePath changes the path of the file.
//
// Parameters:
//   - path: A string representing the new path of the file.
func (fh *fileHandler) ChangePath(path string) {
	fh.loc = path
}
