package internal

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
)

type Agent struct {
	conversation *Conversation
	scanner      *Scanner
}

func NewAgent(scanner *Scanner) *Agent {
	return &Agent{
		conversation: &Conversation{
			ID:       uuid.NewString(),
			Messages: []Message{},
		},
		scanner: scanner,
	}
}

func (a *Agent) Chat(handler func(message Message, conversation *Conversation) (Message, error)) {
	fmt.Print("> ")
	a.scanner.Scan(func(input string) error {
		// User message
		message := Message{
			ID:      uuid.NewString(),
			Role:    RoleUser,
			Content: input,
		}

		// AI response
		response, err := handler(message, a.conversation)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Store the messages in the conversation
		a.conversation.Messages = append(a.conversation.Messages, message, response)

		fmt.Println(RenderMarkdown(response.Content))
		fmt.Print("> ")

		return nil
	})
}

func (a *Agent) SystemPrompt(name string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(name, "~") {
		name = filepath.Join(home, name[1:])
	}

	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		panic("system prompt file does not exist" + name)
	}

	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(data)
}
