package internal

import "github.com/charmbracelet/glamour"

func RenderMarkdown(content string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(180),
	)

	data, err := r.Render(content)
	if err != nil {
		panic(err)
	}
	return data
}
