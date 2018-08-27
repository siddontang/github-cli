package main

import (
	"fmt"

	trending "github.com/andygrunwald/go-trending"
	"github.com/spf13/cobra"
)

var (
	trendingTime string
)

func newTrendingCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "trending [language]",
		Short: "Github CLI for trending popular projects",
		Args:  cobra.MinimumNArgs(0),
		Run:   runTrendingCommandFunc,
	}

	m.Flags().StringVar(&trendingTime, "time", "daily", "Trending time: daily, weekly or monthly")

	return m
}

func formatLanguage(lan string) string {
	if len(lan) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s] ", lan)
}

func runTrendingCommandFunc(cmd *cobra.Command, args []string) {
	lan := ""
	if len(args) == 1 {
		lan = args[0]
	}

	trend := trending.NewTrending()

	projects, err := trend.GetProjects(trendingTime, lan)
	perror(err)

	for _, project := range projects {
		fmt.Printf("%s%s %s\n", formatLanguage(project.Language), project.URL.String(), project.Description)
	}
}
