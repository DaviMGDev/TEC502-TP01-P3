package main

import (
	"cod-client/internal/commands"
	"cod-client/internal/services"
	"cod-client/internal/state"
	"cod-client/internal/utils"
	"strings"
)

func main() {
	// 1. Inicializa o estado (inclui UI e cliente MQTT)
	appState := state.New()

	// 2. Inicializa os serviços, injetando o estado
	eventSvc := services.NewEventService(appState)
	subSvc := services.NewSubscriptionService(appState)

	// 3. Inicializa o gerenciador de comandos
	cmdManager := commands.NewManager(eventSvc, appState)

	// 4. Configura o Mux de comandos (roteador de texto)
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

	// 5. Inicia as inscrições nos tópicos MQTT
	subSvc.SubscribeToAll()

	// 6. Inicia o loop da UI e o processamento de entrada
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

	// Mantém a aplicação rodando indefinidamente
	select {}
}