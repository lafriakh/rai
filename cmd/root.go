package cmd

import (
	"fmt"
	"os"
	"rai/internal"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "ai",
	Short:            "Ai is a command line tool for interacting with AI models",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Rai CLI", strings.Join(args, " "))
	},
}

func Execute(config internal.Config) {
	rootCmd.Context()
	rootCmd.AddCommand(NewGeminiCmd(config))
	rootCmd.AddCommand(NewAnthropicCmd(config))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
