package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Chat struct {
	Outputs chan string
	Inputs  chan string
	Reader  *bufio.Reader
}

func NewChat() *Chat {
	return &Chat{
		Outputs: make(chan string, 10),
		Inputs:  make(chan string, 1),
		Reader:  bufio.NewReader(os.Stdin),
	}
}

func (chat *Chat) Write(message string) {
	chat.Outputs <- message
}

func (chat *Chat) WriteLoop() {
	for output := range chat.Outputs {
		fmt.Println(output)
	}
}

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

func (chat *Chat) Start(fn func()) {
	go chat.WriteLoop()
	go chat.ReadLoop()
	go fn()
}

func (chat *Chat) Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
