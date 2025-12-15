package utils 

import (
	"encoding/json"
)

type Dict map[string]any 

func (dict Dict) String() (string, error) {
	bytes, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (dict Dict) Json() ([]byte, error) {
	return json.Marshal(dict)
}
