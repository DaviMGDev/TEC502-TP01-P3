package cluster

import (
	"cod-server/internal/api"
	"encoding/json"
	"fmt"
	"io"

	raft "github.com/hashicorp/raft"
)

// FSMSnapshot represents a serialized state snapshot for Raft fault tolerance.
// It implements the raft.FSMSnapshot interface for persistence.
type FSMSnapshot struct {
	data []byte
}

// Persist writes the snapshot data to the provided sink and closes it on success.
func (s *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		if _, err := sink.Write(s.data); err != nil {
			return err
		}
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
		return err
	}

	return nil
}

// Release is called when the snapshot is no longer needed; currently a no-op.
func (s *FSMSnapshot) Release() {}

// ClusterFSM is the Raft Finite State Machine that converts committed Raft logs into system actions.
// It delegates business logic to an event handler.
type ClusterFSM struct {
	// eventHandler processes game logic based on event method
	eventHandler api.EventHandlerInterface
}

// NewClusterFSM creates a new ClusterFSM with dependency injection of the event handler.
func NewClusterFSM(handler api.EventHandlerInterface) *ClusterFSM {
	return &ClusterFSM{
		eventHandler: handler,
	}
}

// Apply is called by Raft when a log entry is committed to the consensus log.
// It deserializes the log data and routes to the appropriate event handler method.
func (fsm *ClusterFSM) Apply(log *raft.Log) interface{} {
	var event api.Event
	if err := json.Unmarshal(log.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal log data: %w", err)
	}

	switch event.Method {
	case "register":
		return fsm.eventHandler.OnRegister(event)
	case "login":
		return fsm.eventHandler.OnLogin(event)
	case "get_cards":
		return fsm.eventHandler.OnGetCards(event)
	case "buy_pack":
		return fsm.eventHandler.OnBuyPack(event)
	case "offer_trade":
		return fsm.eventHandler.OnOfferTrade(event)
	case "start_match":
		return fsm.eventHandler.OnStartMatch(event)
	case "join_match":
		return fsm.eventHandler.OnJoinMatch(event)
	case "surrender_match":
		return fsm.eventHandler.OnSurrenderMatch(event)
	case "make_move":
		return fsm.eventHandler.OnMakeMove(event)
	default:
		return fmt.Errorf("unhandled fsm event method: %s", event.Method)
	}
}

// Snapshot returns a point-in-time copy of the current system state.
// TODO: Implement full state serialization of repositories for production use.
func (fsm *ClusterFSM) Snapshot() (raft.FSMSnapshot, error) {
	// Note: Currently returns an empty snapshot placeholder.
	// Full implementation would serialize all repository state.
	snapData := []byte("{}")

	return &FSMSnapshot{data: snapData}, nil
}

// Restore reconstructs the FSM state from a snapshot.
// TODO: Implement full state deserialization and repopulation of repositories.
func (fsm *ClusterFSM) Restore(rc io.ReadCloser) error {
	// Read snapshot data; currently discarded as placeholder implementation.
	// Full implementation would deserialize and repopulate repository state.
	_, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read snapshot data: %w", err)
	}

	return nil
}
