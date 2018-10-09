package main

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	pullsState     string
	pullsLimit     int
	pullsSinceTime string
	pullsOffsetDur string
	pullsOwner     string
	pullsReviewers string
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
	m.Flags().StringVar(&pullsOffsetDur, "offset", "-48h", "The offset of since time")
	m.Flags().StringVar(&pullsOwner, "owner", "", "The Github account")
	m.Flags().StringVar(&pullsReviewers, "reviewers", "", "Request reivewers, separated by comma")
	return m
}

func runPullsCommandFunc(cmd *cobra.Command, args []string) {
	opts := SearchOptions{
		Order: "desc",
		Sort:  "updated",
		Limit: issuesLimit,
	}

	queryArgs := url.Values{}
	users := splitUsers(pullsReviewers)
	for _, user := range users {
		queryArgs.Add("assignee", user)
	}

	queryArgs.Add("is", "pr")
	rangeTime := newRangeTime()
	rangeTime.adjust(pullsSinceTime, pullsOffsetDur)

	queryArgs.Add("updated", rangeTime.String())
	queryArgs.Add("state", pullsState)

	repos := filterRepo(globalClient.cfg, pullsOwner, args)

	m, err := globalClient.SearchIssues(globalCtx, repos, opts, queryArgs)
	perror(err)

	for repo, pulls := range m {
		if len(pulls) == 0 {
			continue
		}

		fmt.Fprintln(&output, repo)
		for _, pull := range pulls {
			fmt.Fprintf(&output, "%s %s %s\n", pull.GetUpdatedAt().Format(TimeFormat), pull.GetHTMLURL(), pull.GetTitle())
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

	fmt.Fprintf(&output, "Title: %s\n", pull.GetTitle())
	fmt.Fprintf(&output, "Created at %s\n", pull.GetCreatedAt().Format(TimeFormat))
	fmt.Fprintf(&output, "Message:\n %s\n", pull.GetBody())
	if len(comments) > pullCommentLimit {
		comments = comments[0:pullCommentLimit]
	}
	for _, comment := range comments {
		fmt.Fprintf(&output, "Comment:\n %s\n", comment.GetBody())
	}
}
