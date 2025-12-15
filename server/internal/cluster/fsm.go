package cluster

import (
	"cod-server/internal/api"
	"encoding/json"
	"fmt"
	"io"

	raft "github.com/hashicorp/raft"
)

// FSMSnapshot implementa a interface raft.FSMSnapshot
type FSMSnapshot struct {
	data []byte
}

// Persist saves the snapshot to the given sink
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

// Release is invoked when we are finished with the snapshot
func (s *FSMSnapshot) Release() {}

// ClusterFSM é a máquina de estados que traduz logs do Raft em ações do sistema
type ClusterFSM struct {
	// Dependência: A interface que sabe lidar com a lógica do jogo
	eventHandler api.EventHandlerInterface
}

// NewClusterFSM cria a FSM injetando o handler de eventos
func NewClusterFSM(handler api.EventHandlerInterface) *ClusterFSM {
	return &ClusterFSM{
		eventHandler: handler,
	}
}

// Apply é chamado pelo Raft quando um log é commitado
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

// Snapshot retorna um "retrato" do estado atual do sistema
func (fsm *ClusterFSM) Snapshot() (raft.FSMSnapshot, error) {
	// Para implementar snapshots completos, precisamos serializar o estado dos repositórios
	// Esta é uma implementação simplificada que retorna um snapshot vazio
	// Em uma implementação completa, serializaríamos o estado dos repositórios

	snapData := []byte("{}") // Placeholder para dados reais do estado

	return &FSMSnapshot{data: snapData}, nil
}

// Restore restaura o estado a partir de um backup
func (fsm *ClusterFSM) Restore(rc io.ReadCloser) error {
	// Ler os dados do snapshot
	_, err := io.ReadAll(rc) // Apenas leitura para cumprir interface, implementação real futuramente
	if err != nil {
		return fmt.Errorf("failed to read snapshot data: %w", err)
	}

	// Em uma implementação completa, desserializaríamos os dados e repopularíamos
	// os repositórios/services com esses dados
	// Por enquanto, implementação vazia como placeholder

	return nil
}
