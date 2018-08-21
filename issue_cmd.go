package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	issuesState      string
	issuesLimit      int
	issuesSinceTime  string
	issuessOffsetDur string
	issuesOwner      string
)

func newIssuesCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "issues [repo]",
		Short: "Github CLI for listing issues",
		Args:  cobra.MinimumNArgs(0),
		Run:   runIssuesCommandFunc,
	}
	m.Flags().StringVar(&issuesState, "state", "open", "Issue state: open or closed")
	m.Flags().IntVar(&issuesLimit, "limit", 20, "Maximum issues limit for a repository")
	m.Flags().StringVar(&issuesSinceTime, "since", "", fmt.Sprintf("Pull Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&issuessOffsetDur, "offset", "-336h", "The offset of since time")
	m.Flags().StringVar(&issuesOwner, "owner", "", "The Github account")
	return m
}

func runIssuesCommandFunc(cmd *cobra.Command, args []string) {
	opts := NewIssueOptions()
	opts.State = issuesState
	opts.Limit = issuesLimit

	if len(issuesSinceTime) > 0 {
		end, err := time.Parse(TimeFormat, issuesSinceTime)
		perror(err)
		opts.End = end
	}

	d, err := time.ParseDuration(issuessOffsetDur)
	perror(err)
	opts.Start = opts.End.Add(d)
	if opts.Start.After(opts.End) {
		opts.Start, opts.End = opts.End, opts.Start
	}

	repos := filterRepo(globalClient.cfg, issuesOwner, args)

	m, err := globalClient.ListIssues(globalCtx, opts, repos)
	perror(err)

	for repo, issues := range m {
		fmt.Println(repo)
		for _, issue := range issues {
			fmt.Printf("%s %s %s\n", issue.GetUpdatedAt().Format(TimeFormat), issue.GetHTMLURL(), issue.GetTitle())
		}
	}
}

var (
	issueCommentLimit int
)

func newIssueCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "issue [repo] [id]",
		Short: "Github CLI for getting one pull",
		Args:  cobra.MinimumNArgs(2),
		Run:   runIssueCommandFunc,
	}

	m.Flags().IntVar(&pullCommentLimit, "comments-limit", 3, "Comments limit")
	return m
}

func runIssueCommandFunc(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[1])
	perror(err)

	repo := findRepo(globalClient.cfg, args)

	issue, err := globalClient.GetIssue(globalCtx, repo.Owner, repo.Name, id)
	perror(err)

	comments, err := globalClient.ListIssueComments(globalCtx, repo.Owner, repo.Name, id)
	perror(err)

	fmt.Printf("Title: %s\n", issue.GetTitle())
	fmt.Printf("Created at %s\n", issue.GetCreatedAt().Format(TimeFormat))
	fmt.Printf("Message:\n %s\n", issue.GetBody())
	if len(comments) > issueCommentLimit {
		comments = comments[0:issueCommentLimit]
	}
	for _, comment := range comments {
		fmt.Printf("Comment:\n %s\n", comment.GetBody())
	}
}
