package Sorting

import (
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// BucketSet is a type that represents a set of buckets.
type BucketSet[K comparable, E any] struct {
	// buckets is the map of buckets.
	buckets map[K]*Bucket[E]
}

// NewBucketSet creates a map of buckets from the given elements using the given
// function.
//
// Parameters:
//   - elems: elements to add to the buckets.
//   - f: function to use to determine the size of the buckets.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
func NewBucketSet[K comparable, E any](elems []E, f func(E) K) *BucketSet[K, E] {
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
