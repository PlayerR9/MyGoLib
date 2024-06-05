package FileManager

import (
	"io/fs"
	"os"
	"path"

	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
)

// Item is a struct that represents a directory entry with its location.
type Item struct {
	// loc is the path to the directory.
	loc string

	// fs.DirEntry is the directory entry.
	fs.DirEntry
}

// NewItem creates a new item.
//
// Parameters:
//   - loc: A string representing the path to the directory.
//   - de: The directory entry.
//
// Returns:
//   - *Item: The new item.
func NewItem(loc string, de fs.DirEntry) Item {
	return Item{loc: loc, DirEntry: de}
}

// Path returns the path to the directory. This
// excludes the directory entry name.
//
// Returns:
//   - string: The path to the directory.
func (i *Item) Path() string {
	return i.loc
}

// FullPath returns the full path to the directory entry including
// the directory entry name.
//
// Returns:
//   - string: The full path to the directory entry.
func (i *Item) FullPath() string {
	return path.Join(i.loc, i.Name())
}

// DEWalkerIter is an iterator that reads directories and all of its subdirectories
// in a depth-first manner without using recursion.
type DEWalkerIter struct {
	// loc is the path to the parent directory.
	loc string

	// source is the slice of directory entries to be iterated.
	source []fs.DirEntry

	// toSee is the slice of directory entries to be visited.
	toSee []Item
}

// Size implements the Iterators.Iterater interface.
func (i *DEWalkerIter) Size() int {
	return len(i.toSee)
}

// Consume implements the Iterators.Iterater interface.
func (i *DEWalkerIter) Consume() (ui.Iterater[Item], error) {
	if len(i.toSee) == 0 {
		return nil, ui.NewErrExhaustedIter()
	}

	var todo []Item

	currentPath := i.toSee[0].FullPath()

	if i.toSee[0].IsDir() {
		subEntries, err := os.ReadDir(currentPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range subEntries {
			item := NewItem(currentPath, entry)

			todo = append(todo, item)
		}

		i.toSee = i.toSee[1:]
	} else {
		firstDirIndex := -1

		for j := 1; j < len(i.toSee); j++ {
			if i.toSee[j].IsDir() {
				firstDirIndex = j
				break
			}
		}

		if firstDirIndex == -1 {
			todo = i.toSee
			i.toSee = i.toSee[:0]
		} else {
			todo = i.toSee[:firstDirIndex]
			i.toSee = i.toSee[firstDirIndex:]
		}
	}

	return ui.NewSimpleIterator(todo), nil
}

// Restart implements the Iterators.Iterater interface.
func (i *DEWalkerIter) Restart() {
	var toSee []Item

	for _, item := range i.source {
		toSee = append(toSee, NewItem(i.loc, item))
	}

	i.toSee = toSee
}

// NewDirEntryIterator creates a new directory entry iterator.
//
// This iterator reads the directories and all of its subdirectories
// in a depth-first manner without using recursion.
//
// Parameters:
//   - loc: A string representing the path to the directory.
//
// Returns:
//   - *Iter1: The new directory entry iterator.
//   - error: An error if it fails to read the directory.
func NewDirEntryIterator(loc string) (*DEWalkerIter, error) {
	entries, err := os.ReadDir(loc)
	if err != nil {
		return nil, err
	}

	var toSee []Item

	for _, entry := range entries {
		toSee = append(toSee, NewItem(loc, entry))
	}

	return &DEWalkerIter{
		loc:    loc,
		source: entries,
		toSee:  toSee,
	}, nil
}
