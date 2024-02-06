package Interfaces

import (
	"fmt"
	"reflect"
)

type EqualityComparable interface {
	Equal(any) bool
	Similar(any) bool
}

type InequalityComparable interface {
	Less(any) bool
}

func wrapType(t any) any {
	switch x := t.(type) {
	case int:
		return ctype[int]{Value: x}
	case int8:
		return ctype[int8]{Value: x}
	case int16:
		return ctype[int16]{Value: x}
	case int32:
		return ctype[int32]{Value: x}
	case int64:
		return ctype[int64]{Value: x}
	case uint:
		return ctype[uint]{Value: x}
	case uint8:
		return ctype[uint8]{Value: x}
	case uint16:
		return ctype[uint16]{Value: x}
	case uint32:
		return ctype[uint32]{Value: x}
	case uint64:
		return ctype[uint64]{Value: x}
	case float32:
		return ctype[float32]{Value: x}
	case float64:
		return ctype[float64]{Value: x}
	case string:
		return ctype[string]{Value: x}
	case complex64:
		return gtype[complex64]{Value: x}
	case complex128:
		return gtype[complex128]{Value: x}
	case bool:
		return gtype[bool]{Value: x}
	default:
		return t
	}
}

// types that are both equality and inequality comparable
type ctype[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 |
		~string
}] struct {
	Value T
}

func (t *ctype[T]) Equal(other any) bool {
	// 1. If the other type is a pointer, dereference it.
	if reflect.ValueOf(other).Kind() == reflect.Ptr {
		if other == nil {
			return t == nil
		}

		if other == nil && t == nil {
			return true
		} else if t != nil {
			r
		}

	dereferenced := reflect.Indirect(reflect.ValueOf(other)).Interface()

	switch y := other.(type) {
	case ctype[T]:
	case interface{ Equal(any) bool }:
	case 

	switch y := other.(type) {
	case interface{ Equal(any) bool }:
		return y.Equal(t.Value)
	default:
		y, ok := other.(T)
		return ok && t.Value == y
	}
}

func (t *ctype[T]) Similar(other any) bool {
	switch y := other.(type) {
	case interface{ Similar(any) bool }:
		return y.Similar(t.Value)
	default:
		y, ok := other.(T)
		return ok && t.Value == y
	}
}

func (t *ctype[T]) Less(other any) bool {
	y, ok := other.(T)
	if !ok {
		panic(fmt.Errorf("type mismatch: %T != %T", t.Value, other))
	}

	return t.Value < y
}

// types that are only equality comparable
type gtype[T interface {
	~complex64 | ~complex128 | ~bool
}] struct {
	Value T
}

func (t *gtype[T]) Equal(other any) bool {
	switch y := other.(type) {
	case interface{ Equal(any) bool }:
		return y.Equal(t.Value)
	default:
		y, ok := other.(T)
		return ok && t.Value == y
	}
}

func (t *gtype[T]) Similar(other any) bool {
	switch y := other.(type) {
	case interface{ Similar(any) bool }:
		return y.Similar(t.Value)
	default:
		y, ok := other.(T)
		return ok && t.Value == y
	}
}

func (t *gtype[T]) Less(other any) bool {
	panic(fmt.Errorf("type %T does not support the Less method", t.Value))
}

type EqualityComparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 |
		~string | ~complex64 | ~complex128 | ~bool

	Equal(any) bool
	Similar(any) bool
}

type InequalityComparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 |
		~string
}

func Equal(left, right interface {
	Equal(any) bool
	any
}) bool {
	switch left.(type) {
	case interface{ Equal(any) bool }:
		return left.Equal(right)
	case int:
		y, ok := right.(int)
	}

	// If two values do not have the same type, they are not equal.
	return left.Equal(right) || left == right
}

func NotEqual[T interface {
	EqualityComparable
	any
}](left, right T) bool {
	// If two values do not have the same type, they are different.
	return !left.Equal(right) && left != right
}

func Similar[T EqualityComparable](left, right T) bool {
	if left.Similar(right) {
		return true
	}

	switch x := t1.(type) {
	case int:
		y, ok := t2.(int)
		if !ok {
			return false
		}
	case int8:
		y, ok := t2.(int8)
		if !ok {
			return false
		}
	case int16:
		y, ok := t2.(int16)
		if !ok {
			return false
		}
	case int32:
		y, ok := t2.(int32)
		if !ok {
			return false
		}
	case int64:
		y, ok := t2.(int64)
		if !ok {
			return false
		}
	case uint:
		y, ok := t2.(uint)
		if !ok {
			return false
		}
	case uint8:
		y, ok := t2.(uint8)
		if !ok {
			return false
		}
	case uint16:
		y, ok := t2.(uint16)
		if !ok {
			return false
		}
	case uint32:
		y, ok := t2.(uint32)
		if !ok {
			return false
		}
	case uint64:
		y, ok := t2.(uint64)
		if !ok {
			return false
		}
	case float32:
		y, ok := t2.(float32)
		if !ok {
			return false
		}
	case float64:
		y, ok := t2.(float64)
		if !ok {
			return false
		}
	case string:
		y, ok := t2.(string)
		if !ok {
			return false
		}
	case complex64:
		y, ok := t2.(complex64)
		if !ok {
			return false
		}
	case complex128:
		y, ok := t2.(complex128)
		if !ok {
			return false
		}
	case bool:
		y, ok := t2.(bool)
		if !ok {
			return false
		}
	}

	return false
}

func Less[T interface {
	InequalityComparable
	Less(T) bool
	any
}](left, right T) bool {
	if left.Less(right) {
		return true
	}

	if right.Less(left) {
		return false
	}

	return left.Less(right) || left < right
}

func LessEqual[T interface {
	InequalityComparable
	Less(T) bool
	any
}](left, right T) bool {
	return !right.Less(left) || left <= right
}

func GreaterThan[T interface {
	InequalityComparable
	Less(T) bool
	any
}](left, right T) bool {
	return right.Less(left) || left > right
}

func GreaterEqual[T interface {
	InequalityComparable
	Less(T) bool
	any
}](left, right T) bool {
	return !left.Less(right) || left >= right
}

type ComparatorOp int

const (
	EQ ComparatorOp = iota
	S
	NE
	LT
	LE
	GT
	GE
)

func Compare[T interface {
	any
	EqualityComparable
	interface{ Equal(any) bool }
	interface{ Similar(any) bool }
	interface{ Less(any) bool }
	InequalityComparable
}](left T, op ComparatorOp, right T) bool {
	switch op {
	case EQ:
		return left.Equal(right) || left == right
	case S:
		return left.Similar(right) || left == right
	case NE:
		return !left.Equal(right) || left != right
	case LT:
		return left.Less(right) || left < right
	case LE:
		return !right.Less(left) || left <= right
	case GT:
		return right.Less(left) || left > right
	case GE:
		return !left.Less(right) || left >= right
	default:
		panic(fmt.Errorf("operator (%d) is not supported", op))
	}
}
