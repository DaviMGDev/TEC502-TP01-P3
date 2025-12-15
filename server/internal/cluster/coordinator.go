package cluster

import (
	"cod-server/internal/api"
	"cod-server/internal/api/mqtt"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	raft "github.com/hashicorp/raft"
)

// CoordinatorInterface define como o sistema lida com eventos de entrada
type CoordinatorInterface interface {
	// Handle recebe um evento do mundo externo e garante que ele seja processado pelo cluster
	Handle(event api.Event) error
}

// RaftCoordinator é a implementação que decide entre aplicar localmente ou encaminhar
type RaftCoordinator struct {
	raftNode    *raft.Raft                // Para verificar estado e aplicar logs
	transport   ClusterTransportInterface // Para encaminhar se não for líder
	mqttAdapter mqtt.MQTTAdapterInterface // Para publicar respostas de volta ao cliente
	timeout     time.Duration             // Tempo máximo de espera pelo consenso
}

// NewRaftCoordinator cria a instância
func NewRaftCoordinator(r *raft.Raft, t ClusterTransportInterface, mqttAdapter mqtt.MQTTAdapterInterface) *RaftCoordinator {
	return &RaftCoordinator{
		raftNode:    r,
		transport:   t,
		mqttAdapter: mqttAdapter,
		timeout:     10 * time.Second, // Exemplo de valor
	}
}

func (c *RaftCoordinator) Handle(event api.Event) error {
	if c.raftNode.State() != raft.Leader {
		leaderAddr := c.raftNode.Leader()
		if leaderAddr == "" {
			return errors.New("não foi possível encontrar o líder do cluster")
		}

		// O endereço do líder no Raft é o endereço do transporte TCP do Raft.
		// Precisamos encontrar o endereço HTTP correspondente.
		// Esta é uma simplificação; um sistema real precisaria de um mapa
		// (raft address -> http address) gerenciado via gossip ou config.
		// Por enquanto, vamos assumir que a porta HTTP é sempre 8080 e o IP é o mesmo.
		host, _, err := net.SplitHostPort(string(leaderAddr))
		if err != nil {
			return fmt.Errorf("endereço do líder inválido: %w", err)
		}
		httpAddr := fmt.Sprintf("%s:8080", host)

		eventBytes, err := event.Json()
		if err != nil {
			return fmt.Errorf("falha ao serializar evento para encaminhamento: %w", err)
		}
		return c.transport.ForwardCommand(httpAddr, eventBytes)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("falha ao serializar evento: %w", err)
	}

	applyFuture := c.raftNode.Apply(data, c.timeout)
	if err := applyFuture.Error(); err != nil {
		return fmt.Errorf("erro ao aplicar comando no raft: %w", err)
	}

	// .Response() pode conter a resposta da FSM.
	// Aqui, estamos assumindo que, se não houver erro, a operação foi bem-sucedida.
	// Um sistema mais complexo poderia retornar dados da FSM para o cliente.
	response := applyFuture.Response()

	// Publicar resposta de volta via MQTT para o cliente receber
	if response != nil {
		if err, ok := response.(error); ok {
			// É um erro, então não é o tipo esperado de resposta
			return err
		} else if responseEvent, ok := response.(api.Event); ok {
			// Obter tópico de resposta apropriado
			replyTopic := c.getReplyTopic(event)
			if replyTopic != "" {
				if err := c.mqttAdapter.Publish(replyTopic, responseEvent); err != nil {
					// Log do erro, mas não retornar erro para não afetar o fluxo principal
					fmt.Printf("Erro ao publicar resposta no MQTT: %v\n", err)
				}
			}
		}
	}

	return nil
}

// getReplyTopic determina o tópico de resposta apropriado com base no método do evento
func (c *RaftCoordinator) getReplyTopic(event api.Event) string {
	switch event.Method {
	case "register":
		return "user/register/events"
	case "login":
		return "user/login/events"
	case "chat":
		// Para chat, responder no próprio tópico de chat
		if roomID, ok := event.Payload["room_id"].(string); ok {
			return "chat/room/" + roomID
		} else {
			return "chat/room/messages"
		}
	default:
		// Para outros métodos, retornar tópico genérico ou baseado no método
		return "responses/" + event.Method
	}
}
