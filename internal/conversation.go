package internal

import (
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"google.golang.org/genai"
)

type Conversation struct {
	Name      *string   `json:"name,omitempty"`
	Messages  []*Message `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Conversation) ToGemini() []*genai.Content {
	var messages []*genai.Content
	for _, m := range c.Messages {
		messages = append(messages, &genai.Content{
			Role: string(m.Role),
			Parts: []*genai.Part{
				{Text: m.Content},
			},
		})
	}
	return messages
}

func (c *Conversation) ToClaude() []anthropic.MessageParam {
	var messages []anthropic.MessageParam
	for _, m := range c.Messages {
		messages = append(messages, anthropic.MessageParam{
			Role: m.Role.ToClaude(),
			Content: []anthropic.ContentBlockParamUnion{
				anthropic.NewTextBlock(m.Content),
			},
		})
	}
	return messages
}

type Message struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
}

func (m *Message) ToGemini() *genai.Part {
	return &genai.Part{
		Text: m.Content,
	}
}
func (m *Message) ToClaude() anthropic.MessageParam {
	return anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content))
}
