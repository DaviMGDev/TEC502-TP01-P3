package state

import (
	"cod-client/internal/ui"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// State detém todas as instâncias e variáveis de estado da aplicação.
type State struct {
	UserID   string
	RoomID   string
	Client   mqtt.Client
	Chat     *ui.Chat
}

// New inicializa e retorna um novo estado da aplicação,
// incluindo a conexão com o broker MQTT.
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