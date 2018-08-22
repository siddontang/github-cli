package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	pullsState     string
	pullsLimit     int
	pullsSinceTime string
	pullsOffsetDur string
	pullsOwner     string
)

func newPullsCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "pulls [repo]",
		Short: "Github CLI for listing pulls",
		Args:  cobra.MinimumNArgs(0),
		Run:   runPullsCommandFunc,
	}
	m.Flags().StringVar(&pullsState, "state", "open", "PR state: open or closed")
	m.Flags().IntVar(&pullsLimit, "limit", 20, "Maximum pull limit for a repository")
	m.Flags().StringVar(&pullsSinceTime, "since", "", fmt.Sprintf("Pull Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&pullsOffsetDur, "offset", "-336h", "The offset of since time")
	m.Flags().StringVar(&pullsOwner, "owner", "", "The Github account")
	return m
}

func runPullsCommandFunc(cmd *cobra.Command, args []string) {
	opts := NewPullOptions()
	opts.State = pullsState
	opts.Limit = pullsLimit

	opts.RangeTime.adjust(pullsSinceTime, pullsOffsetDur)

	repos := filterRepo(globalClient.cfg, pullsOwner, args)

	m, err := globalClient.ListPulls(globalCtx, opts, repos)
	perror(err)

	for repo, pulls := range m {
		fmt.Println(repo)
		for _, pull := range pulls {
			fmt.Printf("%s %s %s\n", pull.GetUpdatedAt().Format(TimeFormat), pull.GetHTMLURL(), pull.GetTitle())
		}
	}
}

var (
	pullCommentLimit int
)

func newPullCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "pull [repo] [id]",
		Short: "Github CLI for getting one pull",
		Args:  cobra.MinimumNArgs(2),
		Run:   runPullCommandFunc,
	}

	m.Flags().IntVar(&pullCommentLimit, "comments-limit", 3, "Comments limit")
	return m
}

func runPullCommandFunc(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[1])
	perror(err)

	repo := findRepo(globalClient.cfg, args)

	pull, err := globalClient.GetPull(globalCtx, repo.Owner, repo.Name, id)
	perror(err)

	comments, err := globalClient.ListPullComments(globalCtx, repo.Owner, repo.Name, id)
	perror(err)

	fmt.Printf("Title: %s\n", pull.GetTitle())
	fmt.Printf("Created at %s\n", pull.GetCreatedAt().Format(TimeFormat))
	fmt.Printf("Message:\n %s\n", pull.GetBody())
	if len(comments) > pullCommentLimit {
		comments = comments[0:pullCommentLimit]
	}
	for _, comment := range comments {
		fmt.Printf("Comment:\n %s\n", comment.GetBody())
	}
}
