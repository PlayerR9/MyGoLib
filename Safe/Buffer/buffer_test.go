package Buffer

import (
	"testing"
)

func TestInit(t *testing.T) {
	buffer, err := NewBuffer[int](0)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	buffer.Start()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	for val := range receiveFrom {
		t.Logf("Received %d", val)
	}

	buffer.Wait()
}

func TestTrimFrom(t *testing.T) {
	buffer, err := NewBuffer[int](1)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	buffer.Start()

	sendTo, receiveFrom := buffer.GetSendChannel(), buffer.GetReceiveChannel()

	sendTo <- 1
	sendTo <- 2
	sendTo <- 3

	close(sendTo)

	x, ok := <-receiveFrom
	if !ok {
		t.Errorf("Expected true, got %t", ok)
	}

	if x != 1 {
		t.Errorf("Expected 1, got %d", x)
	}

	buffer.CleanBuffer()

	buffer.Wait()

	_, ok = <-receiveFrom

	if ok {
		t.Errorf("Expected false, got %t", ok)
	}
}
