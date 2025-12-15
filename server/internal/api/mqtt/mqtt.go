package mqtt

import (
	"cod-server/internal/api"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTAdapterInterface interface {
	Connect() error
	Publish(topic string, event api.Event) error
	Subscribe(topic string, handler mqtt.MessageHandler) error
	Disconnect()
}

type MQTTAdapter struct {
	client mqtt.Client
	logger *log.Logger
}

// NewMQTTAdapter cria uma nova instância do adaptador MQTT
func NewMQTTAdapter(broker, clientID string) (MQTTAdapterInterface, error) {
	if clientID == "" {
		clientID = uuid.New().String()
	}

	logger := log.With("component", "mqtt")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		logger.Infof("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
	})

	opts.OnConnect = func(client mqtt.Client) {
		logger.Info("Connected to broker")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.Warnf("Connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	return &MQTTAdapter{client: client, logger: logger}, nil
}

// Connect estabelece a conexão com o broker MQTT
func (a *MQTTAdapter) Connect() error {
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("mqtt connection error: %w", token.Error())
	}
	return nil
}

// Publish publica um evento em um tópico MQTT
func (a *MQTTAdapter) Publish(topic string, event api.Event) error {
	payload, err := event.Json()
	if err != nil {
		return fmt.Errorf("failed to serialize event to json: %w", err)
	}

	token := a.client.Publish(topic, 1, false, payload)

	// Aguardar o token de forma não bloqueante com timeout ou usar uma abordagem assíncrona mais controlada
	// A chamada original esperava dentro de uma goroutine, mas vamos manter isso para não bloquear
	// a thread principal, mas vamos adicionar mais log para melhor monitoramento
	go func() {
		// Implementando timeout com canal para evitar esperas infinitas
		done := make(chan struct{})
		go func() {
			token.Wait() // Aguarda até que a operação seja concluída
			close(done)
		}()

		select {
		case <-done:
			// Operação concluída
			if token.Error() != nil {
				a.logger.Errorf("Failed to publish to topic %s: %v", topic, token.Error())
			} else {
				a.logger.Debugf("Successfully published to topic %s", topic)
			}
		case <-time.After(10 * time.Second): // Timeout de 10 segundos
			a.logger.Errorf("Publish timeout to topic %s", topic)
		}
	}()

	return nil
}

// Subscribe se inscreve em um tópico MQTT
func (a *MQTTAdapter) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if token := a.client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	a.logger.Infof("Subscribed to topic: %s", topic)
	return nil
}

// Disconnect encerra a conexão com o broker MQTT
func (a *MQTTAdapter) Disconnect() {
	a.logger.Info("Disconnecting...")
	a.client.Disconnect(250)
}
