package common

// Slicer is an interface that provides a method to convert a data structure to a slice.
type Slicer[T any] interface {
	// Slice returns a slice containing all the elements in the data structure.
	//
	// Returns:
	//   - []T: A slice containing all the elements in the data structure.
	Slice() []T

	Iterable[T]
}

// SliceOf converts any type to a slice of elements of the same type.
//
// Parameters:
//   - elem: The element to convert to a slice.
//
// Returns:
//   - []T: The slice representation of the element.
//
// Behaviors:
//   - Nil elements are converted to nil slices.
//   - Slice elements are returned as is.
//   - Slicer elements have their Slice method called.
//   - Other elements are converted to slices containing a single element.
func SliceOf[T any](elem any) []T {
	if elem == nil {
		return nil
	}

	switch elem := elem.(type) {
	case []T:
		return elem
	case Slicer[T]:
		slice := elem.Slice()
		return slice
	default:
		return []T{elem.(T)}
	}
}

// Iterable is an interface that defines a method to get an iterator over a
// collection of elements of type T. It is implemented by data structures that
// can be iterated over.
type Iterable[T any] interface {
	// Iterator returns an iterator over the collection of elements.
	//
	// Returns:
	//   - Iterater[T]: An iterator over the collection of elements.
	Iterator() Iterater[T]
}

// IteratorOf converts any type to an iterator over elements of the same type.
//
// Parameters:
//   - elem: The element to convert to an iterator.
//
// Returns:
//   - Iterater[T]: The iterator over the element.
//
// Behaviors:
//   - IF elem is nil, an empty iterator is returned.
//   - IF elem -implements-> Iterater[T], the element is returned as is.
//   - IF elem -implements-> Iterable[T], the element's Iterator method is called.
//   - IF elem -implements-> []T, a new iterator over the slice is created.
//   - ELSE, a new iterator over a single-element collection is created.
func IteratorOf[T any](elem any) Iterater[T] {
	if elem == nil {
		var builder Builder[T]

		return builder.Build()
	}

	var iter Iterater[T]

	switch elem := elem.(type) {
	case Iterater[T]:
		iter = elem
	case Iterable[T]:
		iter = elem.Iterator()
	case []T:
		iter = &SimpleIterator[T]{
			values: &elem,
			index:  0,
		}
	default:
		iter = &SimpleIterator[T]{
			values: &[]T{elem.(T)},
			index:  0,
		}
	}

	return iter
}

// Iterater is an interface that defines methods for an iterator over a
// collection of elements of type T.
type Iterater[T any] interface {
	// Size returns the number of elements in the collection.
	//
	// Returns:
	//  - count: The number of elements in the collection.
	Size() (count int)

	// Consume advances the iterator to the next element in the
	// collection and returns the current element.
	//
	// Returns:
	//  - T: The current element in the collection.
	//  - error: An error if the iterator is exhausted or if an error occurred
	//    while consuming the element.
	Consume() (T, error)

	// Restart resets the iterator to the beginning of the
	// collection.
	Restart()
}

// ProceduralIterator is a struct that allows iterating over a collection of
// iterators of type Iterater[T].
type ProceduralIterator[E Iterable[T], T any] struct {
	// source is the iterator over the collection of iterators.
	source Iterater[E]

	// iter is the iterator in the collection.
	iter Iterater[T]
}

// Size implements the Iterater interface.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (pi *ProceduralIterator[E, T]) Size() (count int) {
	if pi.iter != nil {
		count = pi.iter.Size()
	}

	count += pi.source.Size()

	return
}

// Consume implements the Iterater interface.
func (pi *ProceduralIterator[E, T]) Consume() (T, error) {
	if pi.iter == nil {
		iter, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = iter.Iterator().(*SimpleIterator[T])
	}

	var val T
	var err error

	for {
		val, err = pi.iter.Consume()
		if err == nil {
			break
		}

		ok := Is[*ErrExhaustedIter](err)
		if !ok {
			return *new(T), err
		}

		iter, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = iter.Iterator().(*SimpleIterator[T])
	}

	return val, nil
}

// Restart implements the Iterater interface.
func (pi *ProceduralIterator[E, T]) Restart() {
	pi.iter = nil
	pi.source.Restart()
}

// IteratorFromIterator creates a new iterator over a collection of iterators
// of type Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and
// return the elements from each iterator in turn.
//
// Parameters:
//   - source: The iterator over the collection of iterators to iterate over.
//
// Return:
//   - *ProceduralIterator[E, T]: The new iterator over the collection of elements.
//     Nil if source is nil.
func NewProceduralIterator[E Iterable[T], T any](source Iterater[E]) *ProceduralIterator[E, T] {
	if source == nil {
		return nil
	}

	pi := &ProceduralIterator[E, T]{
		source: source,
		iter:   nil,
	}

	return pi
}

// SliceIterator is a struct that allows iterating over a collection of
// iterators of type Iterater[T].
type SliceIterator[T any] struct {
	// source is the iterator over the collection of iterators.
	source Iterater[[]T]

	// iter is the iterator in the collection.
	iter *SimpleIterator[T]
}

// Size implements the Iterater interface.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (pi *SliceIterator[T]) Size() (count int) {
	if pi.iter != nil {
		count = pi.iter.Size()
	}

	count += pi.source.Size()

	return
}

// Consume implements the Iterater interface.
func (pi *SliceIterator[T]) Consume() (T, error) {
	if pi.iter == nil {
		values, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = NewSimpleIterator(values)
	}

	var val T
	var err error

	for {
		val, err = pi.iter.Consume()
		if err == nil {
			break
		}

		ok := Is[*ErrExhaustedIter](err)
		if !ok {
			return *new(T), err
		}

		iter, err := pi.source.Consume()
		if err != nil {
			return *new(T), err
		}

		pi.iter = NewSimpleIterator(iter)
	}

	return val, nil
}

// Restart implements the Iterater interface.
func (pi *SliceIterator[T]) Restart() {
	pi.iter = nil
	pi.source.Restart()
}

// IteratorFromIterator creates a new iterator over a collection of iterators
// of type Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and
// return the elements from each iterator in turn.
//
// Parameters:
//   - source: The iterator over the collection of iterators to iterate over.
//
// Return:
//   - *SliceIterator[T]: The new iterator over the collection of elements.
//     Nil if source is nil.
func NewSliceIterator[T any](source Iterater[[]T]) *SliceIterator[T] {
	if source == nil {
		return nil
	}

	pi := &SliceIterator[T]{
		source: source,
		iter:   nil,
	}

	return pi
}

// DynamicIterator is a struct that allows iterating over a collection
// of iterators of type Iterater[T].
type DynamicIterator[E, T any] struct {
	// source is the iterator over the collection of iterators.
	source Iterater[E]

	// iter is the iterator in the collection.
	iter Iterater[T]

	// transition is the transition function that takes an element of type E and
	// returns an iterator.
	transition func(E) Iterater[T]
}

// Size implements the Iterater interface for the DynamicIterator type.
//
// Size is evaluated by summing the sizes of the current iterator and the source
// iterator. Of course, this is just an estimate of the total number of elements
// in the collection.
func (di *DynamicIterator[E, T]) Size() (count int) {
	if di.iter != nil {
		count = di.iter.Size()
	}

	count += di.source.Size()

	return
}

// Consume implements the Iterater interface.
func (di *DynamicIterator[E, T]) Consume() (T, error) {
	if di.iter == nil {
		iter, err := di.source.Consume()
		if err != nil {
			return *new(T), err
		}

		di.iter = di.transition(iter)
	}

	var val T
	var err error

	for {
		val, err = di.iter.Consume()
		if err == nil {
			break
		}

		ok := Is[*ErrExhaustedIter](err)
		if !ok {
			return *new(T), err
		}

		iter, err := di.source.Consume()
		if err != nil {
			return *new(T), err
		}

		di.iter = di.transition(iter)
	}

	return val, nil
}

// Restart implements the Iterater interface.
func (di *DynamicIterator[E, T]) Restart() {
	di.iter = nil
	di.source.Restart()
}

// IteratorFromIterator creates a new iterator over a collection of iterators
// of type Iterater[T].
// It uses the input iterator to iterate over the collection of iterators and
// return the elements from each iterator in turn.
//
// Parameters:
//   - source: The iterator over the collection of iterators to iterate over.
//   - f: The transition function that takes an element of type E and returns
//     an iterator.
//
// Return:
//   - *DynamicIterator[E, T]: The new iterator. Nil if f or source is nil.
func NewDynamicIterator[E, T any](source Iterater[E], f func(E) Iterater[T]) *DynamicIterator[E, T] {
	if f == nil || source == nil {
		return nil
	}

	iter := &DynamicIterator[E, T]{
		source: source,
		iter:   nil,
	}

	iter.transition = f

	return iter
}

// SimpleIterator is a struct that allows iterating over a slice of
// elements of any type.
type SimpleIterator[T any] struct {
	// values is a slice of elements of type T.
	values *[]T

	// index is the current index of the iterator.
	// 0 means not initialized.
	index int
}

// Size implements the Iterater interface.
func (iter *SimpleIterator[T]) Size() int {
	if iter.values == nil {
		return 0
	}

	return len(*iter.values)
}

// Consume implements the Iterater interface.
func (iter *SimpleIterator[T]) Consume() (T, error) {
	if iter.index >= len(*iter.values) {
		return *new(T), NewErrExhaustedIter()
	}

	value := (*iter.values)[iter.index]

	iter.index++

	return value, nil
}

// Restart implements the Iterater interface.
func (iter *SimpleIterator[T]) Restart() {
	iter.index = 0
}

// NewSimpleIterator creates a new iterator over a slice of elements of type T.
//
// Parameters:
//   - values: The slice of elements to iterate over.
//
// Return:
//   - *SimpleIterator[T]: A new iterator over the given slice of elements.
//
// Behaviors:
//   - If values is nil, the iterator is initialized with an empty slice.
//   - Modifications to the slice of elements after creating the iterator will
//     affect the values seen by the iterator.
func NewSimpleIterator[T any](values []T) *SimpleIterator[T] {
	if len(values) == 0 {
		values = make([]T, 0)
	}

	si := &SimpleIterator[T]{
		values: &values,
		index:  0,
	}
	return si
}
