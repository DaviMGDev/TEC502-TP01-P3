package protocol

import (
	"encoding/json"
	"time"
)

// Event é a estrutura comum de mensagens trocadas via MQTT entre projetos.
// Padroniza roteamento pelo método, carimbo de tempo e um mapa de payload flexível.
type Event struct {
	Method    string                 `json:"method"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// Json serializa o Event em um slice de bytes JSON compacto.
func (event Event) Json() ([]byte, error) {
	return json.Marshal(event)
}

// FromJson desserializa um slice de bytes JSON em uma instância de Event.
// Retorna um ponteiro para Event ou um erro se a análise falhar.
func FromJson(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
