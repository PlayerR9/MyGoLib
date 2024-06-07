package Buffer

import (
	"sync"
	"testing"
)

func TestInit(t *testing.T) {
	const (
		MaxCount int = 100
	)

	buffer, err := NewBuffer[int](0)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	sender, receiver := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	var wg sync.WaitGroup

	buffer.Start()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < MaxCount; i++ {
			x, ok := receiver.Receive()
			if !ok {
				t.Errorf("Expected true, got %t", ok)
				return
			}

			if x != i {
				t.Errorf("Expected %d, got %d", i, x)
				return
			}
		}
	}()

	for i := 0; i < MaxCount; i++ {
		sender.Send(i)
	}

	buffer.Close()

	wg.Wait()
}

func TestTrimFrom(t *testing.T) {
	const (
		MaxCount int = 100
	)

	buffer, err := NewBuffer[int](1)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	sender, receiver := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	var wg sync.WaitGroup

	wg.Add(1)

	buffer.Start()

	go func(max int) {
		defer wg.Done()

		for {
			x, ok := receiver.Receive()
			if !ok {
				break
			}

			t.Logf("Received %d", x)
		}
	}(MaxCount)

	for i := 0; i < MaxCount; i++ {
		sender.Send(i)
	}

	buffer.CleanBuffer()

	buffer.Close()

	t.Fatalf("Expected no error, got %s", "error")

	wg.Wait()

	_, ok := receiver.Receive()
	if ok {
		t.Errorf("Expected false, got %t", ok)
	}
}
