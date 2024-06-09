package OrderedMap

import (
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ui "github.com/PlayerR9/MyGoLib/Units/Iterators"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// ModifyValueFunc is a function that modifies a value.
//
// Parameters:
//   - V: The value to modify.
//
// Returns:
//   - V: The modified value.
type ModifyValueFunc[V any] func(V) (V, error)

// MapIterPrinter is a type that represents a printer for a map iterator.
type MapIterPrinter[T ffs.FStringer] struct {
	// Name is the name of the second element in the key-value pair.
	Name string

	// Iter is the iterator to print.
	Iter ui.Iterater[*uc.Pair[string, T]]
}

// FStringIterator is a helper function that iterates over a list of key-value pairs
// and prints them as a list of key-value pairs.
//
// Parameters:
//   - trav: The traversor to use for printing.
//   - elem: The list of key-value pairs to print.
//   - second: The name of the second element in the key-value pair.
//
// Returns:
//   - error: An error if the printing fails.
//
// Behaviors:
//   - The key-value pairs are printed as a list of key-value pairs.

// FString is a function that prints a map iterator.
//
// Format:
//
//	<name>:
//		<key 1>:
//			<value 1>
//		<key 2>:
//			<value 2>
//		// ...
//
// Parameters:
//   - trav: The traversor to use for printing.
//
// Returns:
//   - error: An error if the printing fails.
func (mip *MapIterPrinter[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if mip.Iter == nil {
		return nil
	}

	for at := 0; ; at++ {
		entry, err := mip.Iter.Consume()
		if err != nil {
			break
		}

		err = trav.AddJoinedLine("", "- ", entry.First, ":")
		if err != nil {
			return err
		}

		err = ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithIncreasedIndent(),
			),
			trav,
			entry.Second,
		)
		if err != nil {
			return ue.NewErrAt(at+1, mip.Name, err)
		}
	}

	return nil
}

// NewMapIterPrinter is a function that creates a new map iterator printer.
//
// Parameters:
//   - name: The name of the second element in the key-value pair.
//   - iter: The iterator to print.
//
// Returns:
//   - *MapIterPrinter: A pointer to the newly created map iterator printer.
func NewMapIterPrinter[T ffs.FStringer](name string, iter ui.Iterater[*uc.Pair[string, T]]) *MapIterPrinter[T] {
	return &MapIterPrinter[T]{
		Name: name,
		Iter: iter,
	}
}

// OrderedMapPrinter is a type that represents a printer for an ordered map.
type OrderedMapPrinter[T ffs.FStringer] struct {
	// The name of the ordered map.
	Name string

	// The ordered map to print.
	Map *OrderedMap[string, T]

	// ValueName is the name of the value in the ordered map.
	ValueName string

	// IfEmpty is the string to print if the ordered map is empty.
	IfEmpty string
}

// FString is a function that prints an ordered map.
//
// Format:
//
//	<name>:
//		<key 1>:
//			<value 1>
//		<key 2>:
//			<value 2>
//		// ...
//
// or
//
//	<name>: <ifEmpty>
//
// Parameters:
//   - trav: The traversor to use for printing.
//
// Returns:
//   - error: An error if the printing fails.
func (p *OrderedMapPrinter[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if p.Map == nil {
		return nil
	}

	err := trav.AppendJoinedString("", p.Name, ":")
	if err != nil {
		return err
	}

	if p.Map.Size() == 0 {
		err := trav.AppendRune(' ')
		if err != nil {
			return err
		}

		err = trav.AppendString(p.IfEmpty)
		if err != nil {
			return err
		}

		trav.AcceptLine()
	} else {
		trav.AcceptLine()

		mip := NewMapIterPrinter(p.ValueName, p.Map.Iterator())

		err := ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithIncreasedIndent(),
			),
			trav,
			mip,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewOrderedMapPrinter is a function that creates a new ordered map printer.
//
// Parameters:
//   - name: The name of the ordered map.
//   - elem: The ordered map to print.
//   - valueName: The name of the value in the ordered map.
//   - ifEmpty: The string to print if the ordered map is empty.
//
// Returns:
//   - *OrderedMapPrinter: A pointer to the newly created ordered map printer.
func NewOrderedMapPrinter[T ffs.FStringer](name string, elem *OrderedMap[string, T], valueName string, ifEmpty string) *OrderedMapPrinter[T] {
	return &OrderedMapPrinter[T]{
		Name:      name,
		Map:       elem,
		ValueName: valueName,
		IfEmpty:   ifEmpty,
	}
}
