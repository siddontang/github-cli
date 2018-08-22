package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

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

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", "", "Config File, default ~/.github-cli/config.toml")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "T", "", "Github Token")

	rootCmd.AddCommand(
		newPullsCommand(),
		newPullCommand(),
		newIssuesCommand(),
		newIssueCommand(),
		newTrendingCommand(),
		newEventsCommand(),
	)

	cobra.OnInitialize(initGlobal)
	cobra.EnablePrefixMatching = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
	}
}

func initGlobal() {
	usr, err := user.Current()
	perror(err)

	if len(configFile) == 0 {
		configFile = path.Join(usr.HomeDir, ".github-cli/config.toml")
	}
	cfg, err := NewConfigFromFile(configFile)
	perror(err)

	if len(token) > 0 {
		cfg.Token = token
	}

	if len(cfg.Token) == 0 {
		// try read from ~/.github-cli/token
		name := path.Join(usr.HomeDir, ".github-cli/token")
		if data, err := ioutil.ReadFile(name); err == nil {
			cfg.Token = strings.TrimSpace(string(data))
		}
	}

	globalCtx = context.Background()
	globalClient = NewClient(globalCtx, cfg)
}
