package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Repository is the config for the repository at Github
type Repository struct {
	Owner string `toml:"owner"`
	Name  string `toml:"name"`
}

// Config is the configuration for the tool
type Config struct {
	Account string       `toml:"account"`
	Token   string       `toml:"token"`
	Repos   []Repository `toml:"repos"`
}

// NewConfigFromFile creates the configuration from file
func NewConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	if err = toml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}
