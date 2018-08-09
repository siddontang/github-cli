package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is the client to operate Github
type Client struct {
	c   *github.Client
	cfg *Config
}

// NewClient creates the Github client with token
func NewClient(ctx context.Context, cfg *Config) *Client {
	if len(cfg.Token) == 0 {
		client := github.NewClient(nil)
		return &Client{c: client, cfg: cfg}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &Client{c: client, cfg: cfg}
}
