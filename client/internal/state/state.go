package state

import (
	"cod-client/internal/ui"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// State mantém todas as instâncias e variáveis de estado para toda a aplicação.
// Isso inclui identidade do usuário, contexto da sala, conexão do cliente MQTT e camada UI.
type State struct {
	UserID string      // O identificador único do usuário atualmente logado
	RoomID string      // O identificador da sala de chat/jogo atual
	Client mqtt.Client // Cliente MQTT para operações de publicação/subscrição
	Chat   *ui.Chat    // UI baseada em terminal para interação do usuário
}

// New initializes and returns a new application State instance,
// establishing the connection to the MQTT broker during initialization.
// Panics if the MQTT connection fails to prevent silent failures.
func New() *State {
	chat := ui.NewChat()

	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://broker.emqx.io:8883")
	opts.SetClientID("cod-client-" + fmt.Sprint(time.Now().UnixNano()))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("failed to connect to MQTT broker: %v", token.Error()))
	}

	return &State{
		RoomID: "messages",
		Client: client,
		Chat:   chat,
	}
}
