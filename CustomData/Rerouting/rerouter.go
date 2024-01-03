package Rerouting

import (
	"errors"
	"slices"
	"sync"
)

type ErrAmbiguousConnection struct{}

func (e ErrAmbiguousConnection) Error() string {
	return "cannot have two connections from the same channel"
}

type SRChanneler interface {
	Equals(other SRChanneler) bool
}

type SendChannel struct {
	sendTo chan<- Messager
	code   int
}

func (sc SendChannel) Equals(other SRChanneler) bool {
	otherSC, ok := other.(*SendChannel)
	if !ok {
		return false
	}

	return sc.sendTo == otherSC.sendTo && sc.code == otherSC.code
}

func NewSendChannel(channel chan<- Messager, code int) SendChannel {
	return SendChannel{
		sendTo: channel,
		code:   code,
	}
}

type ReceiveChannel struct {
	receiveFrom <-chan Messager
	code        int
}

func (rc ReceiveChannel) Equals(other SRChanneler) bool {
	otherRC, ok := other.(*ReceiveChannel)
	if !ok {
		return false
	}

	return rc.receiveFrom == otherRC.receiveFrom && rc.code == otherRC.code
}

func NewReceiveChannel(channel <-chan Messager, code int) ReceiveChannel {
	return ReceiveChannel{
		receiveFrom: channel,
		code:        code,
	}
}

type Messager interface {
	Channel() SendChannel
	ParseInexistentEntryPoint() Messager
}

type Rerouter[T any] interface {
	GetSendToChannel() chan<- T
	GetReceiveErrorsChannel() ReceiveChannel
}

type Hub struct {
	connections          map[ReceiveChannel]SendChannel
	entryPoints          []SendChannel
	inexistentEntryPoint SendChannel

	msgBuffer   *Buffer[Messager]
	receiveFrom <-chan Messager
	once        sync.Once
	wg          sync.WaitGroup
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[ReceiveChannel]SendChannel),
		entryPoints: make([]SendChannel, 0),
	}
}

func (h *Hub) AddConnection(from ReceiveChannel, to SendChannel) error {
	_, ok := h.connections[from]
	if ok {
		return &ErrAmbiguousConnection{}
	}

	h.connections[from] = to

	return nil
}

func (h *Hub) AddEntryPoint(from SendChannel) error {
	if slices.ContainsFunc(h.entryPoints, func(x SendChannel) bool {
		return x == from
	}) {
		return errors.New("cannot have two entry points from the same channel")
	}

	h.entryPoints = append(h.entryPoints, from)

	return nil
}

func (h *Hub) SetInexistentEntryPoint(sendTo SendChannel) {
	h.inexistentEntryPoint = sendTo
}

func (h *Hub) Init() chan<- Messager {
	var sendTo chan<- Messager

	h.once.Do(func() {
		sendTo, h.receiveFrom = h.msgBuffer.Init(1)

		h.wg.Add(1)

		go h.routingMessages()
	})

	return sendTo
}

func (h *Hub) routingMessages() {
	defer h.wg.Done()

	// Start the rerouting channel listeners
	var reroutingWg sync.WaitGroup

	for from, to := range h.connections {
		reroutingWg.Add(1)

		go func(from ReceiveChannel, to SendChannel) {
			defer reroutingWg.Done()

			for msg := range from.receiveFrom {
				to.sendTo <- msg
			}
		}(from, to)
	}

	for msg := range h.receiveFrom {
		correctChannel := msg.Channel()

		index := slices.IndexFunc(h.entryPoints, func(x SendChannel) bool {
			return x.Equals(correctChannel)
		})
		if index == -1 {
			h.inexistentEntryPoint.sendTo <- msg.ParseInexistentEntryPoint()
		} else {
			h.entryPoints[index].sendTo <- msg
		}
	}

	reroutingWg.Wait()
}

func (h *Hub) Cleanup() {
	h.wg.Wait()

	var cleanWg sync.WaitGroup
	defer cleanWg.Wait()

	cleanWg.Add(1)

	go func() {
		h.msgBuffer.Cleanup()
		cleanWg.Done()
		h.msgBuffer = nil

		h.receiveFrom = nil
	}()

	allConnections := make(map[SendChannel]bool)

	for _, to := range h.connections {
		_, ok := allConnections[to]
		if !ok {
			allConnections[to] = true
		}
	}

	h.connections = nil

	for _, from := range h.entryPoints {
		_, ok := allConnections[from]
		if !ok {
			allConnections[from] = true
		}
	}

	h.entryPoints = nil

	// Remove inexistent entry point from the map
	delete(allConnections, h.inexistentEntryPoint)

	for to := range allConnections {
		close(to.sendTo)
		to.sendTo = nil
	}

	close(h.inexistentEntryPoint.sendTo)
	h.inexistentEntryPoint.sendTo = nil
}
