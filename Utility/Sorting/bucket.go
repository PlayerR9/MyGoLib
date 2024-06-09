package Sorting

import (
	"slices"

	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Bucket represents a bucket of elements.
type Bucket[T any] struct {
	// elems are the elements of the bucket.
	elems []T
}

// Iterator implements the Iterators.Iterable interface.
func (b *Bucket[T]) Iterator() ui.Iterater[T] {
	return ui.NewSimpleIterator(b.elems)
}

// Copy implements the common.Copier interface.
func (b *Bucket[T]) Copy() uc.Copier {
	elems := make([]T, len(b.elems))
	copy(elems, b.elems)

	return &Bucket[T]{
		elems: elems,
	}
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

// Sort sorts the elements in the bucket using the given comparison function.
//
// Parameters:
//   - f: comparison function to use.
//   - isAsc: flag indicating if the sort is in ascending order.
func (b *Bucket[T]) Sort(sf SortFunc[T], isAsc bool) {
	if sf == nil {
		return
	}

	if !isAsc {
		sf = func(a, b T) int {
			return -sf(a, b)
		}
	}

	slices.SortStableFunc(b.elems, sf)
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
// Behaviors:
//   - If the n is less than or equal to 0, then the bucket will be empty.
//   - If the n is greater then the size of the bucket, then the bucket will
//     remain the same.
func (b *Bucket[T]) Limit(n int) {
	if n <= 0 {
		b.elems = []T{}
	} else if len(b.elems) > n {
		b.elems = b.elems[:n]
	}
}
