package cluster

// DTOs (Data Transfer Objects) used for JSON communication
// -------------------------------------------------

// JoinRequest is the JSON payload sent to /raft/join
type JoinRequest struct {
	NodeID      string `json:"node_id"`
	NodeAddress string `json:"node_address"`
}

// CommandRequest is the JSON payload sent to /raft/command
type CommandRequest struct {
	EventData []byte `json:"event_data"`
}

// Interfaces
// -------------------------------------------------

// ClusterTransportInterface defines how the cluster communicates externally and internally via HTTP
type ClusterTransportInterface interface {
	// Start launches the HTTP server (Gin) in the background
	Start() error

	// JoinCluster is used by a new node to request admission to the cluster
	JoinCluster(targetAddress string, myID string, myAddress string) error

	// ForwardCommand forwards an event to the cluster leader for application
	ForwardCommand(leaderAddress string, eventBytes []byte) error
}
