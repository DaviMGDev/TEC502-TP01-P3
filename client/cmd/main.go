package main

import (
	"cod-client/internal/commands"
	"cod-client/internal/services"
	"cod-client/internal/state"
	"cod-client/internal/utils"
	"strings"
)

func main() {
	// Inicializa o estado da aplicação com UI e conexão do cliente MQTT
	appState := state.New()

	// Cria a camada de serviço para publicação de eventos e tratamento de subscrições
	eventSvc := services.NewEventService(appState)
	subSvc := services.NewSubscriptionService(appState)

	// Configura o gerenciador de comandos para lidar com comandos do usuário com injeção de dependências
	cmdManager := commands.NewManager(eventSvc, appState)

	// Cria e registra todos os comandos disponíveis no roteador de comandos
	mux := utils.NewMux[func(args []string) error](func(args []string) error {
		appState.Chat.Write("Unknown command. Type /help for a list of commands.")
		return nil
	})
	mux.Register("chat", cmdManager.ExecChat)
	mux.Register("login", cmdManager.ExecLogin)
	mux.Register("register", cmdManager.ExecRegister)
	mux.Register("start", cmdManager.ExecStart)
	mux.Register("play", cmdManager.ExecPlay)
	mux.Register("surrender", cmdManager.ExecSurrender)
	mux.Register("join", cmdManager.ExecJoin)
	mux.Register("clear", cmdManager.ExecClear)
	mux.Register("help", cmdManager.ExecHelp)
	mux.Register("exit", cmdManager.ExecExit)

	// Subscreve a todos os tópicos MQTT necessários para receber atualizações e notificações do jogo
	subSvc.SubscribeToAll()

	// Inicia a UI de chat interativa e começa a processar entrada do usuário
	appState.Chat.Start(func() {
		for input := range appState.Chat.Inputs {
			command, args := utils.ParseCommand(strings.TrimSpace(input))
			if handler := mux.Handle(command); handler != nil {
				if err := handler(args); err != nil {
					appState.Chat.Write("Error executing command: " + err.Error())
				}
			}
		}
	})

	// Bloqueia indefinidamente para manter a aplicação em execução
	select {}
}
