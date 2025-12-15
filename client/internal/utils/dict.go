package utils

import (
	"encoding/json"
)

// Dict é um alias de tipo para um mapa de chaves string para quaisquer valores, fornecendo serialização JSON.
type Dict map[string]any

// String serializa o dicionário para uma string JSON formatada com indentação.
func (dict Dict) String() (string, error) {
	bytes, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Json serializa o dicionário para um slice de bytes JSON compacto.
func (dict Dict) Json() ([]byte, error) {
	return json.Marshal(dict)
}
