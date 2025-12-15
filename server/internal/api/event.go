package api

import (
	shared_protocol "shared/protocol"
)

// Tipo wrapper para permitir métodos customizados
type Event struct {
	shared_protocol.Event
}

// Funções wrapper para manter a mesma interface
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
