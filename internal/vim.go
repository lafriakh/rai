package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func EditStringInVim(initialText string) (string, error) {
	tmpFile, err := os.CreateTemp("", "rai-vim-edit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempFilePath := tmpFile.Name()
	defer os.Remove(tempFilePath)

	if _, err := tmpFile.WriteString(initialText); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write initial text to temporary file '%s': %w", tempFilePath, err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file '%s' after writing: %w", tempFilePath, err)
	}

	cmd := exec.Command("vim", "+startinsert", tempFilePath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("vim command finished with error: %w", err)
	}

	modifiedContentBytes, err := os.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read content from temporary file '%s' after editing: %w", tempFilePath, err)
	}

	return string(modifiedContentBytes), nil
}
