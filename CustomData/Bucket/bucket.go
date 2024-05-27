package Bucket

import (
	ui "github.com/PlayerR9/MyGoLib/Units/Iterator"
	us "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// Bucket represents a bucket of elements.
type Bucket[T any] struct {
	// elems are the elements of the bucket.
	elems []T
}

// Iterator implements the Iterator method of the Iterable interface.
func (b *Bucket[T]) Iterator() ui.Iterater[T] {
	return ui.NewSimpleIterator(b.elems)
}

// NewBucket creates a new bucket of elements.
//
// Parameters:
//   - elements: elements to add to the bucket.
//
// Returns:
//   - *Bucket: the new bucket.
func NewBucket[T any](elements []T) *Bucket[T] {
	return &Bucket[T]{
		elems: elements,
	}
}

// Add adds a element to the bucket.
//
// Parameters:
//   - element: element to add.
func (b *Bucket[T]) Add(element T) {
	if len(b.elems) == 0 {
		b.elems = []T{element}
	} else {
		b.elems = append(b.elems, element)
	}
}

// InsertionSort sorts the elements in the bucket using the given comparison
// function.
//
// Parameters:
//   - f: comparison function to use.
func (b *Bucket[T]) InsertionSort(f SortFunc[T]) {
	for i := 1; i < len(b.elems); i++ {
		j := i

		for j > 0 && f(b.elems[j], b.elems[j-1]) {
			b.elems[j], b.elems[j-1] = b.elems[j-1], b.elems[j]
			j--
		}
	}
}

// LinearKeep keeps the elements in the bucket that satisfy the given predicate
// filter.
//
// Parameters:
//   - f: predicate filter to use.
func (b *Bucket[T]) LinearKeep(f us.PredicateFilter[T]) {
	b.elems = us.SliceFilter(b.elems, f)
}

// GetSize returns the size of the bucket.
//
// Returns:
//   - int: the size of the bucket.
func (b *Bucket[T]) GetSize() int {
	return len(b.elems)
}

// Limit limits the number of elements in the bucket.
//
// Parameters:
//   - n: number of elements to keep.
//
// Returns:
//   - *Bucket: the bucket.
func (b *Bucket[T]) Limit(n int) *Bucket[T] {
	if len(b.elems) <= n {
		return b
	} else {
		return &Bucket[T]{
			elems: b.elems[:n],
		}
	}
}
