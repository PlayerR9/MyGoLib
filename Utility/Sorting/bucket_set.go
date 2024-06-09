package Sorting

import (
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// BucketSet is a type that represents a set of buckets.
type BucketSet[K comparable, E any] struct {
	// buckets is the map of buckets.
	buckets map[K]*Bucket[E]
}

// Copy implements the common.Copier interface.
func (bs *BucketSet[K, E]) Copy() uc.Copier {
	buckets := make(map[K]*Bucket[E])

	for size, bucket := range bs.buckets {
		buckets[size] = bucket.Copy().(*Bucket[E])
	}

	return &BucketSet[K, E]{
		buckets: buckets,
	}
}

// Iterator implements the Iterators.Iterable interface.
func (bs *BucketSet[K, E]) Iterator() ui.Iterater[E] {
	var builder ui.Builder[K]

	for size := range bs.buckets {
		builder.Add(size)
	}

	di, err := ui.NewDynamicIterator(
		builder.Build(),
		func(size K) *ui.SimpleIterator[E] {
			return bs.buckets[size].Iterator().(*ui.SimpleIterator[E])
		},
	)
	if err != nil {
		return nil
	}

	return di
}

// NewBucketSet creates a new bucket set from the given map of elements.
//
// Parameters:
//   - m: map of elements to add to the buckets.
//
// Returns:
//   - *BucketSet: the new bucket set.
func NewBucketSet[K comparable, E any](m map[K][]E) *BucketSet[K, E] {
	buckets := make(map[K]*Bucket[E])

	for size, elems := range m {
		buckets[size] = NewBucket(elems)
	}

	return &BucketSet[K, E]{
		buckets: buckets,
	}
}

// MakeBucketSet creates a map of buckets from the given elements using the given
// function.
//
// Parameters:
//   - elems: elements to add to the buckets.
//   - f: function to use to determine the size of the buckets.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
func MakeBucketSet[K comparable, E any](elems []E, f func(E) K) *BucketSet[K, E] {
	buckets := make(map[K]*Bucket[E])

	for _, elem := range elems {
		size := f(elem)

		val, ok := buckets[size]
		if !ok {
			buckets[size] = NewBucket([]E{elem})
		} else {
			val.Add(elem)
		}
	}

	return &BucketSet[K, E]{
		buckets: buckets,
	}
}

// KeepIfBuckets keeps the elements in the buckets that satisfy the given
// predicate filter.
//
// Parameters:
//   - f: predicate filter to use.
func (bs *BucketSet[K, E]) KeepIfBuckets(f us.PredicateFilter[E]) {
	var todo []K // buckets to remove.

	for size, bucket := range bs.buckets {
		bucket.LinearKeep(f)

		if bucket.GetSize() == 0 {
			todo = append(todo, size)
		}
	}

	for _, size := range todo {
		delete(bs.buckets, size)
	}
}

// ViewTopNBuckets returns the top N buckets.
//
// Parameters:
//   - n: number of buckets to return.
func (bs *BucketSet[K, E]) ViewTopNBuckets(n int) {
	if n <= 0 {
		bs.buckets = make(map[K]*Bucket[E])
		return
	}

	for size, bucket := range bs.buckets {
		bucket.Limit(n)
		bs.buckets[size] = bucket
	}
}

// Sort sorts the buckets using the given function.
//
// Parameters:
//   - buckets: buckets to sort.
//   - sf: comparison function to use.
//   - isAsc: flag indicating if the sort is in ascending order.
func (bs *BucketSet[K, E]) Sort(sf SortFunc[E], isAsc bool) {
	if sf == nil {
		return
	}

	if !isAsc {
		sf = func(a, b E) int {
			return -sf(a, b)
		}
	}

	for _, bucket := range bs.buckets {
		bucket.Sort(sf, isAsc)
	}
}

// GetBuckets returns the map of buckets.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
func (bs *BucketSet[K, E]) GetBuckets() map[K]*Bucket[E] {
	return bs.buckets
}

// GetBucket returns the bucket with the given size.
//
// Parameters:
//   - size: size of the bucket to return.
//
// Returns:
//   - *Bucket: the bucket with the given size.
//
// Behaviors:
//   - If the bucket does not exist, nil is returned.
func (bs *BucketSet[K, E]) GetBucket(size K) *Bucket[E] {
	b, ok := bs.buckets[size]
	if !ok {
		return nil
	}

	return b
}

// DoBuckets applies the given function to each bucket.
//
// Parameters:
//   - f: function to apply to each bucket.
//
// Returns:
//   - error: the first error encountered.
func (bs *BucketSet[K, E]) DoBuckets(f func(K, *Bucket[E]) error) error {
	for size, bucket := range bs.buckets {
		err := f(size, bucket)
		if err != nil {
			return err
		}
	}

	return nil
}
