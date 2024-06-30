package FileManager

import (
	"io/fs"
	"os"
	"path"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
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
//   - Item: The new item.
func NewItem(loc string, de fs.DirEntry) Item {
	item := Item{
		loc:      loc,
		DirEntry: de,
	}

	return item
}

// Path returns the path to the directory. This
// excludes the directory entry name.
//
// Returns:
//   - string: The path to the directory.
func (item *Item) Path() string {
	return item.loc
}

// FullPath returns the full path to the directory entry including
// the directory entry name.
//
// Returns:
//   - string: The full path to the directory entry.
func (item *Item) FullPath() string {
	name := item.Name()
	joined := path.Join(item.loc, name)

	return joined
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

	// el is the error list.
	el *uc.ErrOrSol[error]
}

// Size implements the Iterators.Iterater interface.
func (iter *DEWalkerIter) Size() int {
	return len(iter.toSee)
}

// Consume implements the Iterators.Iterater interface.
//
// It ignores entries that would cause an error when reading them.
// However, if all entries are invalid, only the furthest error is returned.
func (iter *DEWalkerIter) Consume() (*ItemList, error) {
	if len(iter.toSee) == 0 {
		return nil, uc.NewErrExhaustedIter()
	}

	for len(iter.toSee) > 0 {
		ok := iter.toSee[0].IsDir()
		if !ok {
			break
		}

		currentPath := iter.toSee[0].FullPath()

		subEntries, err := os.ReadDir(currentPath)
		iter.toSee = iter.toSee[1:]

		if err != nil {
			depth := CountDepth(currentPath)

			iter.el.AddErr(err, depth)
		} else {
			var tmp []Item

			for _, entry := range subEntries {
				item := NewItem(currentPath, entry)
				tmp = append(tmp, item)
			}

			iter.toSee = append(tmp, iter.toSee...)
		}
	}

	if len(iter.toSee) == 0 {
		return nil, iter.el.GetErrors()[0]
	}

	var todo []Item

	firstDirIndex := -1

	for j := 1; j < len(iter.toSee); j++ {
		ok := iter.toSee[j].IsDir()
		if ok {
			firstDirIndex = j
			break
		}
	}

	if firstDirIndex == -1 {
		todo = iter.toSee
		iter.toSee = iter.toSee[:0]
	} else {
		todo = iter.toSee[:firstDirIndex]
		iter.toSee = iter.toSee[firstDirIndex:]
	}

	il := &ItemList{
		items: todo,
	}

	return il, nil
}

// Restart implements the Iterators.Iterater interface.
func (iter *DEWalkerIter) Restart() {
	var toSee []Item

	for _, item := range iter.source {
		tmp := NewItem(iter.loc, item)
		toSee = append(toSee, tmp)
	}

	iter.toSee = toSee

	var el uc.ErrOrSol[error]

	iter.el = &el
}

// ItemList is a struct that represents a list of items.
type ItemList struct {
	// items is the slice of items.
	items []Item
}

// Iterator implements the Iterators.Iterable interface.
func (il *ItemList) Iterator() uc.Iterater[Item] {
	iter := uc.NewSimpleIterator(il.items)

	return iter
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
func NewDEWalkerIter(loc string) (uc.Iterater[Item], error) {
	entries, err := os.ReadDir(loc)
	if err != nil {
		return nil, err
	}

	var toSee []Item

	for _, entry := range entries {
		item := NewItem(loc, entry)

		toSee = append(toSee, item)
	}

	w := &DEWalkerIter{
		loc:    loc,
		source: entries,
		toSee:  toSee,
	}

	iter := uc.NewProceduralIterator(w)

	return iter, nil
}
