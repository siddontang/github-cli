package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

//SendToSlack sends msg to slack
func SendToSlack(cfg Slack, msg string) {
	token := cfg.Token
	channelName := cfg.Channel
	user := cfg.User

	if token == "" {
		perror(fmt.Errorf("must provide a token"))
		return
	}

	if channelName == "" {
		perror(fmt.Errorf("must provide a channel name"))
		return
	}

	if channelName[0] != '#' {
		channelName = "#" + channelName
	}

	api := slack.New(token)
	params := slack.PostMessageParameters{Username: user, Markdown: true}
	_, _, err := api.PostMessage(channelName, msg, params)
	if err != nil {
		perror(fmt.Errorf("can not post msg to slack with err: %v\n", err))
	}
}
