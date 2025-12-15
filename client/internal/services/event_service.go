package services

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/state"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// EventService encapsula a lógica de criação e publicação de eventos.
type EventService struct {
	appState *state.State
}

// NewEventService cria uma nova instância de EventService.
func NewEventService(s *state.State) *EventService {
	return &EventService{appState: s}
}

// createEvent é um helper genérico que constrói um Event com campos padrão.
func (s *EventService) createEvent(method string, payload map[string]interface{}) protocol.Event {
	return protocol.Event{
		Method:    method,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// inferTopicFor determina o tópico MQTT apropriado para publicar um dado evento.
// As rotas são baseadas no método do evento e contexto (userID, roomID, etc).
func (s *EventService) inferTopicFor(event protocol.Event) string {
	switch event.Method {
	case "register", "login":
		return "user/" + event.Method
	case "chat":
		return "chat/room/" + s.appState.RoomID
	case "start":
		return "game/start_game"
	case "play":
		return "game/" + s.appState.RoomID + "/play_card"
	case "surrender":
		return "game/" + s.appState.RoomID + "/surrender"
	case "join":
		return "game/join_game"
	case "buy":
		return "store/buy"
	case "exchange":
		return "cards/" + s.appState.RoomID + "/exchange" + s.appState.UserID
	default:
		return ""
	}
}

// Publish serializa um evento em JSON e o publica no tópico MQTT apropriado.
// Retorna um erro se o tópico for desconhecido ou se a publicação falhar.
func (s *EventService) Publish(event protocol.Event) error {
	topic := s.inferTopicFor(event)
	if topic == "" {
		return fmt.Errorf("unknown topic for method: %s", event.Method)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	token := s.appState.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

// --- Funções de criação de eventos para comandos específicos do usuário ---

// CreateChatEvent constrói um evento de chat a partir dos argumentos do usuário com validação básica.
func (s *EventService) CreateChatEvent(args []string) protocol.Event {
	content := strings.Join(args, " ")

	// Garante que a mensagem não está vazia
	if len(content) < 1 {
		return s.createEvent("chat", map[string]interface{}{
			"error":   "message content cannot be empty",
			"content": content,
			"user_id": s.appState.UserID,
		})
	}

	return s.createEvent("chat", map[string]interface{}{
		"content": content,
		"user_id": s.appState.UserID,
	})
}

// CreateRegisterEvent constrói um evento de registro com validação para usuário e senha.
func (s *EventService) CreateRegisterEvent(args []string) protocol.Event {
	// Garante que usuário e senha foram fornecidos
	if len(args) < 2 {
		return s.createEvent("register", map[string]interface{}{
			"error":    "username and password are required",
			"username": "",
			"password": "",
		})
	}

	username := args[0]
	password := args[1]

	// Validate username is not empty
	if len(username) < 1 {
		return s.createEvent("register", map[string]interface{}{
			"error":    "username cannot be empty",
			"username": username,
			"password": "",
		})
	}

	// Validate password is not empty
	if len(password) < 1 {
		return s.createEvent("register", map[string]interface{}{
			"error":    "password cannot be empty",
			"username": username,
			"password": "",
		})
	}

	return s.createEvent("register", map[string]interface{}{
		"username": username,
		"password": password,
	})
}

// CreateLoginEvent builds a login event with validation for username and password.
func (s *EventService) CreateLoginEvent(args []string) protocol.Event {
	// Ensure both username and password are provided
	if len(args) < 2 {
		return s.createEvent("login", map[string]interface{}{
			"error":    "username and password are required",
			"username": "",
			"password": "",
		})
	}

	username := args[0]
	password := args[1]

	// Validate username is not empty
	if len(username) < 1 {
		return s.createEvent("login", map[string]interface{}{
			"error":    "username cannot be empty",
			"username": username,
			"password": "",
		})
	}

	// Validate password is not empty
	if len(password) < 1 {
		return s.createEvent("login", map[string]interface{}{
			"error":    "password cannot be empty",
			"username": username,
			"password": "",
		})
	}

	return s.createEvent("login", map[string]interface{}{
		"username": username,
		"password": password,
	})
}

// CreateStartGameEvent builds an event to initiate a new game session.
func (s *EventService) CreateStartGameEvent() protocol.Event {
	return s.createEvent("start", map[string]interface{}{
		"user_id": s.appState.UserID,
		//		"room_id": s.appState.RoomID,
	})
}

// CreatePlayCardEvent builds an event to play a specific card in the current game.
func (s *EventService) CreatePlayCardEvent(cardID string) protocol.Event {
	return s.createEvent("play", map[string]interface{}{
		"user_id": s.appState.UserID,
		"room_id": s.appState.RoomID,
		"card_id": cardID,
	})
}

// CreateSurrenderEvent builds an event to surrender the current game.
func (s *EventService) CreateSurrenderEvent() protocol.Event {
	return s.createEvent("surrender", map[string]interface{}{
		"user_id": s.appState.UserID,
		"room_id": s.appState.RoomID,
	})
}

// CreateJoinGameEvent builds an event to join an existing game by room ID.
func (s *EventService) CreateJoinGameEvent(roomID string) protocol.Event {
	return s.createEvent("join", map[string]interface{}{
		"user_id": s.appState.UserID,
		"room_id": roomID,
	})
}

// CreateBuyEvent builds an event to purchase an item from the store.
func (s *EventService) CreateBuyEvent(itemID string) protocol.Event {
	return s.createEvent("buy", map[string]interface{}{
		"user_id": s.appState.UserID,
		"item_id": itemID,
	})
}

// CreateExchangeEvent builds an event to exchange cards with other players.
func (s *EventService) CreateExchangeEvent(cardIDs []string) protocol.Event {
	return s.createEvent("exchange", map[string]interface{}{
		"user_id":  s.appState.UserID,
		"room_id":  s.appState.RoomID,
		"card_ids": cardIDs,
	})
}
