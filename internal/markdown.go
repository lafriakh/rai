package internal

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(content string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(180),
	)

	data, err := r.Render(content)
	if err != nil {
		fmt.Printf("Failed to render markdown: %v\n", err)
	}
	return data
}
