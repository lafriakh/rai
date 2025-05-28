package main

import (
	"rai/cmd"
)

func main() {
	config := initConfig()

	cmd.Execute(config)
}
