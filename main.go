package main

import (
	"ai/cmd"
)

func main() {
	config := initConfig()

	cmd.Execute(config)
}
