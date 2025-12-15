package services

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/state"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SubscriptionService encapsula a lógica de inscrição e manipulação de eventos recebidos.
type SubscriptionService struct {
	appState *state.State
}

// NewSubscriptionService cria uma nova instância de SubscriptionService.
func NewSubscriptionService(s *state.State) *SubscriptionService {
	return &SubscriptionService{appState: s}
}

// subscribe é um helper para se inscrever em um tópico com um handler.
func (s *SubscriptionService) subscribe(topic string, handler mqtt.MessageHandler) {
	if token := s.appState.Client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("failed to subscribe to topic %s: %v", topic, token.Error()))
	}
}

// SubscribeToAll gerencia todas as inscrições necessárias para a aplicação.
func (s *SubscriptionService) SubscribeToAll() {
	s.subscribe("chat/room/"+s.appState.RoomID, s.onChatEvent)
	s.subscribe("user/register/events", s.onRegisterEvent)
	s.subscribe("user/login/events", s.onLoginEvent)
}

// decodeEvent é um helper para deserializar o payload da mensagem.
func (s *SubscriptionService) decodeEvent(msg mqtt.Message) (protocol.Event, error) {
	var event protocol.Event
	err := json.Unmarshal(msg.Payload(), &event)
	if err != nil {
		s.appState.Chat.Write(fmt.Sprintf("Error decoding event: %v", err))
	}
	return event, err
}

// --- Handlers de eventos recebidos ---

func (s *SubscriptionService) onLoginEvent(c mqtt.Client, m mqtt.Message) {
	event, err := s.decodeEvent(m)
	if err != nil {
		return
	}

	if status, ok := event.Payload["status"].(string); ok {
		if status == "success" {
			if newUserID, ok := event.Payload["user_id"].(string); ok {
				s.appState.UserID = newUserID
				s.appState.Chat.Write("Login successful!")
			}
		} else {
			if errorMsg, ok := event.Payload["error"].(string); ok {
				s.appState.Chat.Write("Login failed: " + errorMsg)
			}
		}
	}
}

func (s *SubscriptionService) onRegisterEvent(c mqtt.Client, m mqtt.Message) {
	event, err := s.decodeEvent(m)
	if err != nil {
		return
	}

	if status, ok := event.Payload["status"].(string); ok {
		if status == "success" {
			// A successful registration in this design doesn't automatically log the user in.
			// It just confirms account creation.
			s.appState.Chat.Write("Registration successful! You can now log in.")
		} else {
			if errorMsg, ok := event.Payload["error"].(string); ok {
				s.appState.Chat.Write("Registration failed: " + errorMsg)
			}
		}
	}
}

func (s *SubscriptionService) onChatEvent(c mqtt.Client, m mqtt.Message) {
	event, err := s.decodeEvent(m)
	if err != nil {
		return
	}

	if content, ok := event.Payload["content"].(string); ok {
		// Evita que o usuário veja a própria mensagem que enviou duplicada,
		// caso o broker a envie de volta (QoS 0 não deveria, mas é uma boa proteção).
		if senderID, ok := event.Payload["user_id"].(string); ok {
			if s.appState.UserID != "" && senderID == s.appState.UserID {
				return
			}
		}
		s.appState.Chat.Write(content)
	}
}
