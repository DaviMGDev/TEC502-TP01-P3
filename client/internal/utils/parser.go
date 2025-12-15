package utils

import "strings"

// ParseCommand interpreta a entrada do usuário, separando-a em um nome de comando e argumentos.
// Se a entrada não começar com '/', ela é tratada como um comando 'chat'.
func ParseCommand(input string) (string, []string) {
	if len(input) == 0 {
		return "", []string{}
	}

	if input[0] != '/' {
		return "chat", []string{input}
	}
	
	parts := strings.Fields(input)
	command := strings.ToLower(parts[0][1:])
	args := parts[1:]
	return command, args
}
