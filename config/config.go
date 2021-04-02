package config

import (
	"github.com/spf13/viper"
)

// New initializes the configuration setting. It searches for the
// config file with a filename of "config.yaml".
//
// Config file will be search in the following path in order:
// - "/etc/noteapp"
// - "$HOME/.noteapp"
// - "."
func New() *Config {
	conf, err := newConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func newConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/noteapp")
	viper.AddConfigPath("$HOME/.noteapp")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if viper.Get("store.file.path") == nil {
		viper.Set("store.file.path", ".")
	}

	if viper.Get("store.file.path") == nil {

	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Config is the application-level configuration containing
// all the information for running the application.
type Config struct {
	// Store Database Configuration
	Store Store
}

type Store struct {
	File File
}

type File struct {
	Path string
}
