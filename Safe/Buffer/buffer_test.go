package Buffer

import (
	"testing"
)

func TestInit(t *testing.T) {
	buffer, err := NewBuffer[int](1)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	buffer.Start()
	defer buffer.Wait()

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	for i := 1; i <= 3; i++ {
		x := <-receiveFrom

		if x != i {
			t.Errorf("Expected %d, got %d", i, x)
		}
	}
}

func TestTrimFrom(t *testing.T) {
	buffer, err := NewBuffer[int](1)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	buffer.Start()
	defer buffer.Wait()

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	buffer.CleanBuffer()

	x, ok := <-receiveFrom

	if !ok {
		t.Errorf("Expected true, got %t", ok)
	}

	if x != 1 {
		t.Errorf("Expected 1, got %d", x)
	}

	_, ok = <-receiveFrom

	if ok {
		t.Errorf("Expected false, got %t", ok)
	}
}
