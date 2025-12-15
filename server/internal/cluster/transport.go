package cluster

// DTOs (Data Transfer Objects) para comunicação JSON
// -------------------------------------------------

// JoinRequest representa o payload JSON enviado para /raft/join
type JoinRequest struct {
	NodeID      string `json:"node_id"`
	NodeAddress string `json:"node_address"`
}

// CommandRequest representa o payload JSON enviado para /raft/command
type CommandRequest struct {
	EventData []byte `json:"event_data"`
}

// Interfaces
// -------------------------------------------------

// ClusterTransportInterface define como o cluster se comunica externamente e internamente via HTTP
type ClusterTransportInterface interface {
	// Start inicia o servidor HTTP (Gin) em background
	Start() error

	// JoinCluster é usado por um nó novo para pedir entrada no cluster
	JoinCluster(targetAddress string, myID string, myAddress string) error

	// ForwardCommand é usado por um nó seguidor para repassar um evento ao líder
	ForwardCommand(leaderAddress string, eventBytes []byte) error
}
