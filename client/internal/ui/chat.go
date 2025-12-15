package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Chat representa uma UI de chat baseada em terminal com buffer de mensagens.
// Gerencia canais separados para mensagens de saída e leitura de entrada.
type Chat struct {
	Outputs chan string   // Canal com buffer para mensagens a serem exibidas
	Inputs  chan string   // Canal com buffer para entrada do usuário
	Reader  *bufio.Reader // Reader para buffer de stdin
}

// NewChat creates a new Chat instance with initialized channels.
func NewChat() *Chat {
	return &Chat{
		Outputs: make(chan string, 10),
		Inputs:  make(chan string, 1),
		Reader:  bufio.NewReader(os.Stdin),
	}
}

// Write envia uma mensagem para o canal de saída para exibição.
func (chat *Chat) Write(message string) {
	chat.Outputs <- message
}

// WriteLoop continuamente lê do canal de saída e imprime mensagens.
// Deve ser executado em uma goroutine separada.
func (chat *Chat) WriteLoop() {
	for output := range chat.Outputs {
		fmt.Println(output)
	}
}

// ReadLoop continuamente lê entrada do usuário do stdin e a envia via canal de entrada.
// Deve ser executado em uma goroutine separada e sai graciosamente em EOF.
func (chat *Chat) ReadLoop() {
	for {
		input, err := chat.Reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		chat.Inputs <- input
	}
}

// Start lança a UI com goroutines separadas para saída, entrada e uma função customizada.
// A função fornecida fn será executada em sua própria goroutine.
func (chat *Chat) Start(fn func()) {
	go chat.WriteLoop()
	go chat.ReadLoop()
	go fn()
}

// Clear executa o comando shell 'clear' para apagar a exibição do terminal.
func (chat *Chat) Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
