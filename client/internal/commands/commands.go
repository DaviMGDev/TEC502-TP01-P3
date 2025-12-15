package commands

import (
	"cod-client/internal/services"
	"cod-client/internal/state"
	// "fmt"
	"os"
	"time"
)

// Manager encapsula a lógica de execução para todos os comandos do usuário.
type Manager struct {
	eventSvc *services.EventService
	appState *state.State
}

// NewManager cria uma nova instância do Command Manager.
func NewManager(eventSvc *services.EventService, appState *state.State) *Manager {
	return &Manager{eventSvc: eventSvc, appState: appState}
}

// --- Implementações dos Comandos ---

func (m *Manager) ExecChat(args []string) error {
	if m.appState.UserID == "" {
		m.appState.Chat.Write("You must be logged in to chat. Use /login <user> <pass>")
		return nil // Retornar nil para não mostrar "Error: ..." na UI
	}
	event := m.eventSvc.CreateChatEvent(args)
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecLogin(args []string) error {
	if len(args) < 2 {
		m.appState.Chat.Write("Usage: /login <username> <password>")
		return nil
	}
	event := m.eventSvc.CreateLoginEvent(args)
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecRegister(args []string) error {
	if len(args) < 2 {
		m.appState.Chat.Write("Usage: /register <username> <password>")
		return nil
	}
	event := m.eventSvc.CreateRegisterEvent(args)
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecClear(args []string) error {
	m.appState.Chat.Clear()
	return nil
}

func (m *Manager) ExecHelp(args []string) error {
	helpText := `Available commands:
/register <username> <password> - Register a new user
/login <username> <password>    - Login as an existing user
/chat <message>               - Send a chat message (or just type without a '/')
/start                        - Start a new game
/play <card_id>               - Play a card in game
/surrender                    - Surrender current game
/join <game_id>               - Join a game
/clear                        - Clear the chat window
/help                         - Show this help message
/exit                         - Exit the application`
	m.appState.Chat.Write(helpText)
	return nil
}

func (m *Manager) ExecStart(args []string) error {
	if m.appState.UserID == "" {
		m.appState.Chat.Write("You must be logged in to start a game. Use /login <user> <pass>")
		return nil
	}
	event := m.eventSvc.CreateStartGameEvent()
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecPlay(args []string) error {
	if m.appState.UserID == "" {
		m.appState.Chat.Write("You must be logged in to play a card. Use /login <user> <pass>")
		return nil
	}
	if len(args) < 1 {
		m.appState.Chat.Write("Usage: /play <card_id>")
		return nil
	}
	event := m.eventSvc.CreatePlayCardEvent(args[0])
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecSurrender(args []string) error {
	if m.appState.UserID == "" {
		m.appState.Chat.Write("You must be logged in to surrender. Use /login <user> <pass>")
		return nil
	}
	event := m.eventSvc.CreateSurrenderEvent()
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecJoin(args []string) error {
	if m.appState.UserID == "" {
		m.appState.Chat.Write("You must be logged in to join a game. Use /login <user> <pass>")
		return nil
	}
	if len(args) < 1 {
		m.appState.Chat.Write("Usage: /join <game_id>")
		return nil
	}
	event := m.eventSvc.CreateJoinGameEvent(args[0])
	return m.eventSvc.Publish(event)
}

func (m *Manager) ExecExit(args []string) error {
	m.appState.Chat.Write("Exiting chat...")
	time.Sleep(1 * time.Second)
	m.appState.Client.Disconnect(250)
	os.Exit(0)
	return nil // Unreachable, but satisfies the interface
}
