package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	issuesState     string
	issuesLimit     int
	issuesSinceTime string
	issuesOffsetDur string
	issuesOwner     string
	issuesAssignees string
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
	m.Flags().StringVar(&issuesSinceTime, "since", "", fmt.Sprintf("Issue Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&issuesOffsetDur, "offset", "-48h", "The offset of since time")
	m.Flags().StringVar(&issuesOwner, "owner", "", "The Github account")
	m.Flags().StringVar(&issuesAssignees, "assignees", "", "Assignees for the issue, separated by comma")
	return m
}

func runIssuesCommandFunc(cmd *cobra.Command, args []string) {
	opts := NewIssueOptions()
	opts.State = issuesState
	opts.Limit = issuesLimit
	opts.Assignees = splitUsers(issuesAssignees)

	opts.RangeTime.adjust(issuesSinceTime, issuesOffsetDur)

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
