package Sorting

import (
	"slices"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// CopyBucketSet copies the given map of buckets.
//
// Parameters:
//   - bs: map of buckets to copy.
//
// Returns:
//   - map[int]*Bucket: the copied map of buckets.
//
// Behaviors:
//   - If the given map is nil, an empty map is returned.
func CopyBucketSet[K comparable, E any](bs map[K][]E) map[K][]E {
	if bs == nil {
		return make(map[K][]E)
	}

	buckets := make(map[K][]E)

	for size, bucket := range bs {
		bucketCopy := make([]E, len(bucket))
		copy(bucketCopy, bucket)

		buckets[size] = bucketCopy
	}

	return buckets
}

// IteratorOfBucketSet is a function type that creates an iterator for the given bucket.
//
// Parameters:
//   - bs: the bucket set.
//
// Returns:
//   - common.Iterater[E]: the iterator.
//
// Behaviors:
//   - If the given bucket set is nil, an empty iterator is returned.
func IteratorOfBucketSet[K comparable, E any](bs map[K][]E) uc.Iterater[E] {
	if bs == nil {
		return uc.NewSimpleIterator([]E{})
	}

	var builder uc.Builder[K]

	for size := range bs {
		builder.Add(size)
	}

	src := builder.Build()

	f := func(size K) uc.Iterater[E] {
		iter := uc.NewSimpleIterator(bs[size])
		return iter
	}

	di := uc.NewDynamicIterator(src, f)

	return di
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
//
// Behaviors:
//   - If the given slice of elements is empty or the function is nil, an empty map
//     is returned.
func MakeBucketSet[K comparable, E any](elems []E, f func(E) K) map[K][]E {
	if len(elems) == 0 || f == nil {
		return make(map[K][]E)
	}

	bs := make(map[K][]E)

	for _, elem := range elems {
		size := f(elem)

		val, ok := bs[size]
		if !ok {
			val = []E{elem}
		} else {
			val = append(val, elem)
		}

		bs[size] = val
	}

	return bs
}

// KeepIfBuckets keeps the elements in the buckets that satisfy the given
// predicate filter.
//
// Parameters:
//   - bs: map of buckets to filter.
//   - f: predicate filter to use.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
//
// Behaviors:
//   - If a bucket is empty after filtering, it is removed.
//   - WARNING: This can modify the original map.
//   - If the given map is nil, an empty map is returned.
//   - If the given filter is nil, the original map is returned (no copy).
func KeepIfBuckets[K comparable, E any](bs map[K][]E, f us.PredicateFilter[E]) map[K][]E {
	if bs == nil {
		return make(map[K][]E)
	}

	if f == nil {
		return bs
	}

	otherBs := make(map[K][]E)

	for key, buckets := range bs {
		buckets := us.SliceFilter(buckets, f)
		if len(buckets) > 0 {
			otherBs[key] = buckets
		}
	}

	return otherBs
}

// ViewTopNBuckets returns the top N buckets.
//
// Parameters:
//   - bs: map of buckets to filter.
//   - n: number of buckets to return.
//
// Returns:
//   - map[int]*Bucket: the map of buckets.
//
// Behaviors:
//   - If the given map is nil or the number is less than or equal to 0, an empty map
//     is returned.
func ViewTopNBuckets[K comparable, E any](bs map[K][]E, n int) map[K][]E {
	if bs == nil || n <= 0 {
		return make(map[K][]E)
	}

	otherBs := make(map[K][]E)

	for size, buckets := range bs {
		if len(buckets) > n {
			buckets = buckets[:n]
		}

		otherBs[size] = buckets
	}

	return otherBs
}

// SortBucketSet sorts the buckets using the given function in-place.
//
// Parameters:
//   - bs: map of buckets to sort.
//   - sf: comparison function to use.
//   - isAsc: flag indicating if the sort is in ascending order.
//
// Behaviors:
//   - If the given map is nil or the function is nil, nothing is done.
func SortBucketSet[K comparable, E any](bs map[K][]E, sf SortFunc[E], isAsc bool) {
	if bs == nil || sf == nil {
		return
	}

	if !isAsc {
		sf = func(a, b E) int {
			return -sf(a, b)
		}
	}

	for _, bucket := range bs {
		slices.SortStableFunc(bucket, sf)
	}
}

// GetBucket returns the bucket with the given key.
//
// Parameters:
//   - bs: map of buckets.
//   - key: key of the bucket to return.
//
// Returns:
//   - *Bucket: the bucket with the given key.
//
// Behaviors:
//   - If the bucket does not exist, nil is returned.
//   - If the given map is nil, nil is returned.
func GetBucket[K comparable, E any](bs map[K][]E, key K) []E {
	if bs == nil {
		return nil
	}

	b, ok := bs[key]
	if !ok {
		return nil
	}

	return b
}

// DoBuckets applies the given function to each bucket.
//
// Parameters:
//   - bs: map of buckets.
//   - f: function to apply to each bucket.
//
// Returns:
//   - error: the first error encountered.
//
// Behaviors:
//   - If the given map is nil or the function is nil, nothing is done.
func DoBuckets[K comparable, E any](bs map[K][]E, f func(K, []E) error) error {
	if bs == nil || f == nil {
		return nil
	}

	for size, bucket := range bs {
		err := f(size, bucket)
		if err != nil {
			return err
		}
	}

	return nil
}
