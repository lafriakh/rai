package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Agent struct {
	conversation *Conversation
	storage      *Storage
	scanner      *Scanner
}

func NewAgent(scanner *Scanner, conversationName string) *Agent {
	var storage *Storage
	if conversationName != "" {
		s, err := NewStorage(fmt.Sprintf("%s.chat", conversationName))
		if err != nil {
			panic(fmt.Errorf("failed to create storage: %w", err))
		}
		storage = s
	}

	agent := &Agent{
		conversation: &Conversation{
			Name:     &conversationName,
			Messages: []Message{},
		},
		scanner: scanner,
		storage: storage,
	}

	return agent
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

		fmt.Println(RenderMarkdown(addPrefixToEachLine(message.Content, "> ")))

		// AI response
		response, err := handler(message, a.conversation)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Store the messages in the conversation
		a.conversation.Messages = append(a.conversation.Messages, message, response)
		if a.storage != nil {
			if err := a.storage.AddMessage(message); err != nil {
				panic(err)
			}
			if err := a.storage.AddMessage(response); err != nil {
				panic(err)
			}
		}

		fmt.Println(RenderMarkdown(response.Content))
		fmt.Print("> ")

		return nil
	})
	
	if err := a.storage.Close(); err != nil {
		panic(err)
	}
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

func addPrefixToEachLine(input string, prefix string) string {
	if input == "" {
		return ""
	}
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
