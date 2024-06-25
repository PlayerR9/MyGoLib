package Sorting

import (
	"slices"

	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// BucketSet is a type that represents a set of buckets.
type BucketSet[K comparable, E any] struct {
	// buckets is the map of buckets.
	buckets map[K][]E
}

// Copy implements the common.Copier interface.
func (bs *BucketSet[K, E]) Copy() uc.Copier {
	buckets := make(map[K][]E)

	for size, bucket := range bs.buckets {
		bucketCopy := make([]E, len(bucket))
		copy(bucketCopy, bucket)

		buckets[size] = bucketCopy
	}

	bsCopy := &BucketSet[K, E]{
		buckets: buckets,
	}

	return bsCopy
}

// Iterator implements the Iterators.Iterable interface.
func (bs *BucketSet[K, E]) Iterator() ui.Iterater[E] {
	var builder ui.Builder[K]

	for size := range bs.buckets {
		builder.Add(size)
	}

	di := ui.NewDynamicIterator(
		builder.Build(),
		func(size K) ui.Iterater[E] {
			iter := ui.NewSimpleIterator(bs.buckets[size])
			return iter
		},
	)

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
	buckets := make(map[K][]E)

	for size, elems := range m {
		bucket := make([]E, len(elems))
		copy(bucket, elems)

		buckets[size] = bucket
	}

	bs := &BucketSet[K, E]{
		buckets: buckets,
	}

	return bs
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
	buckets := make(map[K][]E)

	for _, elem := range elems {
		size := f(elem)

		val, ok := buckets[size]
		if !ok {
			val = []E{elem}
		} else {
			val = append(val, elem)
		}

		buckets[size] = val
	}

	bs := &BucketSet[K, E]{
		buckets: buckets,
	}

	return bs
}

// KeepIfBuckets keeps the elements in the buckets that satisfy the given
// predicate filter.
//
// Parameters:
//   - f: predicate filter to use.
func (bs *BucketSet[K, E]) KeepIfBuckets(f us.PredicateFilter[E]) {
	var todo []K // buckets to remove.

	for key, bucket := range bs.buckets {
		bucket := us.SliceFilter(bucket, f)
		if len(bucket) == 0 {
			todo = append(todo, key)
		}
	}

	for _, key := range todo {
		delete(bs.buckets, key)
	}
}

// ViewTopNBuckets returns the top N buckets.
//
// Parameters:
//   - n: number of buckets to return.
func (bs *BucketSet[K, E]) ViewTopNBuckets(n int) {
	if n <= 0 {
		bs.buckets = make(map[K][]E)
		return
	}

	for size, bucket := range bs.buckets {
		if len(bucket) > n {
			bucket = bucket[:n]
		}

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
		slices.SortStableFunc(bucket, sf)
	}
}

// GetBuckets returns the map of buckets.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
func (bs *BucketSet[K, E]) GetBuckets() map[K][]E {
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
func (bs *BucketSet[K, E]) GetBucket(size K) []E {
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
func (bs *BucketSet[K, E]) DoBuckets(f func(K, []E) error) error {
	for size, bucket := range bs.buckets {
		err := f(size, bucket)
		if err != nil {
			return err
		}
	}

	return nil
}
