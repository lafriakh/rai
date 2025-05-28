package cmd

import (
	"context"
	"fmt"
	"os"
	"rai/internal"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewAnthropicCmd(config internal.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "anthropic",
		Short: "Interact with the Anthropic AI models",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Anthropic (%s) \n", config.Anthropic.ModelID)

			ctx := context.Background()
			client := newAnthropicClient(cmd, config)

			scanner := internal.NewScanner(os.Stdin)
			agent, err := internal.NewAgent(scanner, cmd.Flag("conversation").Value.String())
			if err != nil {
				return err
			}

			systemPrompt, err := agent.SystemPrompt(cmd.Flag("system").Value.String())
			if err != nil {
				return err
			}
			errr := agent.Chat(func(message *internal.Message, conversation *internal.Conversation) (*internal.Message, error) {
				response, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model:     anthropic.Model(cmd.Flag("model").Value.String()),
					Messages:  append([]anthropic.MessageParam{}, append(conversation.ToClaude(), message.ToClaude())...),
					MaxTokens: int64(8192),
					System: []anthropic.TextBlockParam{
						{Text: systemPrompt},
					},
				})
				if err != nil {
					return nil, err
				}

				return &internal.Message{
					ID:      uuid.NewString(),
					Role:    internal.RoleModel,
					Content: anthropicMessageToText(response.Content),
				}, nil
			})
			
			return errr
		},
	}

	cmd.Flags().String("model", config.Anthropic.ModelID, "Model to use (e.g., Claude3.7 Sonnet)")
	cmd.Flags().String("key", "", "API key for the AI provider")
	cmd.Flags().String("system", config.Anthropic.SystemPromptPath, "Path to the system prompt file to use")
	cmd.Flags().String("conversation", "", "conversation name to store the chat and load messages from")

	return cmd
}

func newAnthropicClient(cmd *cobra.Command, config internal.Config) anthropic.Client {
	key := cmd.Flag("key").Value.String()
	if key == "" {
		key = config.Anthropic.APIKey
	}

	return anthropic.NewClient(option.WithAPIKey(key))
}

func anthropicMessageToText(content []anthropic.ContentBlockUnion) string {
	var res strings.Builder
	for _, content := range content {
		switch content.Type {
		case "text":
			res.WriteString(content.Text)
		}
	}

	return res.String()
}
