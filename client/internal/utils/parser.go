package utils

import "strings"

// ParseCommand analisa a entrada do usuário em um nome de comando e argumentos.
// Entrada começando com '/' é tratada como comando; caso contrário, usa 'chat' como padrão.
// O nome do comando é convertido para minúsculas para correspondência sem diferenciação de maiúsculas.

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
