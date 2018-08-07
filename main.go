package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	token        string
	configFile   string
	globalCtx    context.Context
	globalClient *Client
)

func perror(err error) {
	if err == nil {
		return
	}

	println(err.Error())
	os.Exit(1)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "github-cli",
		Short: "Github CLI",
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", "./config.toml", "Config File")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "T", "", "Github Token")

	rootCmd.AddCommand(
		newPullsCommand(),
		newPullCommand(),
	)

	cobra.OnInitialize(initGlobal)
	cobra.EnablePrefixMatching = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
	}
}

func initGlobal() {
	cfg, err := NewConfigFromFile(configFile)
	perror(err)

	if len(token) > 0 {
		cfg.Token = token
	}

	if len(cfg.Token) == 0 {
		perror(fmt.Errorf("must provide a Github Token"))
	}

	globalCtx = context.Background()
	globalClient = NewClient(globalCtx, cfg)
}
