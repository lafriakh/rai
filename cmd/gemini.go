package cmd

import (
	"ai/internal"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/genai"
)

func NewGeminiCmd(config internal.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gemini",
		Short: "Interact with the Gemini AI models",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Gemini (%s)\n", config.Gemini.ModelID)

			ctx := context.Background()
			client, err := newGeminiClient(ctx, cmd, config)
			if err != nil {
				panic(err)
			}

			scanner := internal.NewScanner(os.Stdin)
			agent := internal.NewAgent(scanner)
			agent.Chat(func(message internal.Message, conversation *internal.Conversation) (internal.Message, error) {
				config := generateContentConfig(agent.SystemPrompt(cmd.Flag("system").Value.String()))
				chat, err := client.Chats.Create(ctx, cmd.Flag("model").Value.String(), config, conversation.ToGemini())
				if err != nil {
					return internal.Message{}, err
				}

				response, err := chat.SendMessage(ctx, *message.ToGemini())
				if err != nil {
					return internal.Message{}, err
				}

				return internal.Message{
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

func generateContentConfig(systemInstruction string) *genai.GenerateContentConfig {
	if viper.GetString("gemini.system") == "" {
		return nil
	}

	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemInstruction},
			},
		},
	}
}
