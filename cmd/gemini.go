package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"rai/internal"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func NewGeminiCmd(config internal.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gemini",
		Short: "Interact with the Gemini AI models",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Gemini (%s)\n", config.Gemini.ModelID)

			ctx := context.Background()
			client, err := newGeminiClient(ctx, cmd, config)
			if err != nil {
				return err
			}

			scanner := internal.NewScanner(os.Stdin)
			agent, err := internal.NewAgent(scanner, cmd.Flag("conversation").Value.String())
			if err != nil {
				return err
			}

			return agent.Chat(func(message *internal.Message, conversation *internal.Conversation) (*internal.Message, error) {
				config, err := generateContentConfig(agent, cmd.Flag("system").Value.String())
				if err != nil {
					return nil, err
				}

				chat, err := client.Chats.Create(ctx, cmd.Flag("model").Value.String(), config, conversation.ToGemini())
				if err != nil {
					return nil, err
				}

				response, err := chat.SendMessage(ctx, *message.ToGemini())
				if err != nil {
					return nil, err
				}

				return &internal.Message{
					ID:      uuid.NewString(),
					Role:    internal.RoleModel,
					Content: response.Text(),
				}, nil
			})
		},
	}

	cmd.Flags().String("model", config.Gemini.ModelID, "Model to use (e.g., gemini-2.5-pro-preview-05-06)")
	cmd.Flags().String("key", "", "API key for the AI provider")
	cmd.Flags().String("system", config.Gemini.SystemPromptPath, "Path to the system prompt file to use")
	cmd.Flags().String("conversation", "", "conversation name to store the chat and load messages from")
	cmd.Flags().Lookup("key").NoOptDefVal = ""

	return cmd
}

func newGeminiClient(ctx context.Context, cmd *cobra.Command, config internal.Config) (*genai.Client, error) {
	key := cmd.Flag("key").Value.String()
	if key == "" {
		key = config.Gemini.APIKey
	}

	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  key,
		Backend: genai.BackendGeminiAPI,
	})
}

func generateContentConfig(agent *internal.Agent, systemPromptPath string) (*genai.GenerateContentConfig, error) {
	if systemPromptPath == "" {
		return nil, errors.New("no system prompt provided")
	}
	systemInstruction, err := agent.SystemPrompt(systemPromptPath)
	if err != nil {
		return nil, err
	}

	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemInstruction},
			},
		},
	}, nil
}
