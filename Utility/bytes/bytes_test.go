package bytes

import (
	"bytes"
	"testing"
)

func TestFindContentIndexes(t *testing.T) {
	var (
		OpToken []byte   = []byte("(")
		ClToken []byte   = []byte(")")
		Content [][]byte = bytes.Fields([]byte("( ( a + b ) * c ) + d"))
	)

	indices, err := FindContentIndexes([]byte(OpToken), []byte(ClToken), Content)
	if err != nil {
		t.Errorf("expected no error, got %s instead", err.Error())
	}

	if indices[0] != 1 {
		t.Errorf("expected 1, got %d instead", indices[0])
	}

	if indices[1] != 9 {
		t.Errorf("expected 9, got %d instead", indices[1])
	}
}
