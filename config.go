package main

import (
	"ai/internal"
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "$HOME/.config/ai/config.yaml"
)

//go:embed config.yaml
var defaultConfig string

var config internal.Config

func initConfig() internal.Config {
	ensureConfig()
	
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	v.AddConfigPath("$HOME/.config/ai")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	v.Unmarshal(&config)

	return config
}

func ensureConfig() {
	filepath := os.ExpandEnv(defaultConfigPath)

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		if err = os.WriteFile(filepath, []byte(defaultConfig), os.FileMode(0o644)); err != nil {
			fmt.Println("Can't write config:", err)
			os.Exit(1)
		}
	}
}
