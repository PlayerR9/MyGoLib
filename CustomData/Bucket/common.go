package Bucket

import (
	slext "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// SortFunc is a function type that defines a comparison function for branches.
//
// Parameters:
//   - b1: first branch to compare.
//   - b2: second branch to compare.
//
// Returns:
//   - bool: true if the first branch is less than the second branch, false
//     otherwise.
type SortFunc[T any] func(b1, b2 T) bool

// MakeBuckets creates a map of buckets from the given elements using the given
// function.
//
// Parameters:
//   - elems: elements to add to the buckets.
//   - f: function to use to determine the size of the buckets.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
func MakeBuckets[T any](elems []T, f func(T) int) map[int]*Bucket[T] {
	buckets := make(map[int]*Bucket[T])

	for _, elem := range elems {
		size := f(elem)

		val, ok := buckets[size]
		if !ok {
			buckets[size] = NewBucket[T]([]T{elem})
		} else {
			val.Add(elem)
		}
	}

	return buckets
}

// SortBuckets sorts the buckets using the given function.
//
// Parameters:
//   - buckets: buckets to sort.
//   - f: comparison function to use.
func SortBuckets[T any](buckets map[int]*Bucket[T], f SortFunc[T]) {
	for _, bucket := range buckets {
		bucket.InsertionSort(f)
	}
}

// KeepIfBuckets keeps the elements in the buckets that satisfy the given
// predicate filter.
//
// Parameters:
//   - buckets: buckets to filter.
//   - f: predicate filter to use.
//
// Returns:
//   - map[int]*Bucket: the filtered buckets.
func KeepIfBuckets[T any](buckets map[int]*Bucket[T], f slext.PredicateFilter[T]) map[int]*Bucket[T] {
	todo := make([]int, 0) // buckets to remove.

	for size, bucket := range buckets {
		bucket.LinearKeep(f)

		if bucket.GetSize() == 0 {
			todo = append(todo, size)
		}
	}

	if len(todo) == 0 {
		return buckets
	}

	for _, size := range todo {
		delete(buckets, size)
	}

	return buckets
}

// ViewTopNBuckets returns the top N buckets.
//
// Parameters:
//   - buckets: buckets to filter.
//   - n: number of buckets to return.
//
// Returns:
//   - map[int]*Bucket: the top N buckets.
func ViewTopNBuckets[T any](buckets map[int]*Bucket[T], n int) map[int]*Bucket[T] {
	if n <= 0 {
		return make(map[int]*Bucket[T])
	}

	newBucket := make(map[int]*Bucket[T])

	for size, bucket := range buckets {
		newBucket[size] = bucket.Limit(n)
	}

	return newBucket
}
