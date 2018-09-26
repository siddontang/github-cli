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

// Slack configuration
type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
	User    string `toml:"user"`
}

// Config is the configuration for the tool
type Config struct {
	Account string       `toml:"account"`
	Token   string       `toml:"token"`
	Slack   Slack        `toml:"slack"`
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

// FindRepo finds a repository.
func (c *Config) FindRepo(owner string, name string) *Repository {
	for _, repo := range c.Repos {
		if len(owner) == 0 && repo.Name == name {
			return &repo
		}

		if repo.Name == name && repo.Owner == owner {
			return &repo
		}
	}

	return nil
}
