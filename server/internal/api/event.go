package api

import (
	shared_protocol "shared/protocol"
)

// Wrapper type to allow custom methods
type Event struct {
	shared_protocol.Event
}

// Wrapper functions to maintain the same interface
func (event Event) Json() ([]byte, error) {
	return event.Event.Json()
}

func FromJson(data []byte) (*Event, error) {
	ev, err := shared_protocol.FromJson(data)
	if err != nil {
		return nil, err
	}
	converted := Event{Event: *ev}
	return &converted, nil
}
