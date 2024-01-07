package Queue

// QueueOperationType is an integer type that represents the type of operation performed
// on a queue.
// It is used in error handling to specify the operation that caused an error.
type QueueOperationType int

const (
	// Dequeue represents a dequeue operation, which removes an element from the queue.
	Dequeue QueueOperationType = iota

	// Peek represents a peek operation, which retrieves the element at the front of the
	// queue without removing it.
	Peek
)

// String is a method of the QueueOperationType type. It returns a string representation
// of the queue operation type.
//
// The method uses an array of strings where the index corresponds to the integer value
// of the QueueOperationType.
// The string at the corresponding index is returned as the string representation of the
// QueueOperationType.
//
// This method is typically used for error messages and logging.
func (qot QueueOperationType) String() string {
	return [...]string{
		"dequeue",
		"peek",
	}[qot]
}

// ErrEmptyQueue is a struct that represents an error when attempting to perform a queue
// operation on an empty queue.
// It has a single field, operation, of type QueueOperationType, which indicates the type
// of operation that caused the error.
type ErrEmptyQueue struct {
	operation QueueOperationType
}

// Error is a method of the ErrEmptyQueue type that implements the error interface. It
// returns a string representation of the error.
// The method constructs the error message by concatenating the string "could not ", the
// string representation of the operation that caused the error,
// and the string ": queue is empty". This provides a clear and descriptive error message
// when attempting to perform a queue operation on an empty queue.
func (e ErrEmptyQueue) Error() string {
	return "could not " + e.operation.String() + ": queue is empty"
}

// ErrFullQueue is a struct that represents an error when attempting to enqueue an element
// into a full queue.
// It does not have any fields as the error condition is solely based on the state of the
// queue being full.
type ErrFullQueue struct{}

// Error is a method of the ErrFullQueue type that implements the error interface. It
// returns a string representation of the error.
// The method returns the string "could not enqueue: queue is full", providing a clear and
// descriptive error message when attempting to enqueue an element into a full queue.
func (e ErrFullQueue) Error() string {
	return "could not enqueue: queue is full"
}

// ErrNegativeCapacity is a struct that represents an error when a negative capacity is
// provided for a queue.
// It does not have any fields as the error condition is solely based on the provided
// capacity being negative.
type ErrNegativeCapacity struct{}

// Error is a method of the ErrNegativeCapacity type that implements the error interface.
// It returns a string representation of the error.
// The method returns the string "capacity of a queue cannot be negative", providing a
// clear and descriptive error message when a negative capacity is provided for a queue.
func (e ErrNegativeCapacity) Error() string {
	return "capacity of a queue cannot be negative"
}

// ErrTooManyValues is a struct that represents an error when too many values are
// provided for initializing a queue.
// It does not have any fields as the error condition is solely based on the number of
// provided values exceeding the capacity of the queue.
type ErrTooManyValues struct{}

// Error is a method of the ErrTooManyValues type that implements the error interface.
// It returns a string representation of the error.
// The method returns the string "could not initialize queue: too many values", providing
// a clear and descriptive error message when too many values are provided for initializing
// a queue.
func (e ErrTooManyValues) Error() string {
	return "could not initialize queue: too many values"
}

// ErrOutOfBoundsIterator is a struct that represents an error when an iterator goes
// out of bounds.
// It does not have any fields as the error condition is solely based on the iterator
// exceeding the bounds of the data structure it is iterating over.
type ErrOutOfBoundsIterator struct{}

// Error is a method of the ErrOutOfBoundsIterator type that implements the error
// interface. It returns a string representation of the error.
// The method returns the string "iterator out of bounds", providing a clear and
// descriptive error message when an iterator goes out of bounds.
func (e ErrOutOfBoundsIterator) Error() string {
	return "iterator out of bounds"
}
