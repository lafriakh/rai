package main

import (
	"ai/internal"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultConfigPath = "$HOME/.config/rai/config.yaml"
)

//go:embed config.yaml
var defaultConfig string

var config internal.Config

func initConfig() internal.Config {
	ensureConfig()
	
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	v.AddConfigPath("$HOME/.config/rai")
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
	fpath := os.ExpandEnv(defaultConfigPath)

	if _, err := os.Stat(fpath); errors.Is(err, os.ErrNotExist) {
		makeDirectoryIfNotExists(filepath.Dir(fpath))
		if err = os.WriteFile(fpath, []byte(defaultConfig), os.FileMode(0o644)); err != nil {
			fmt.Println("Can't write config:", err)
			os.Exit(1)
		}
	}
}
