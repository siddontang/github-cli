package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// IssueOptions is the options for listing issues
type IssueOptions struct {
	RangeTime

	// open, closed, all
	// default open
	State string
	// created, updated, comments
	// Default: updated
	Sort string
	// maximum number
	Limit int

	Assignees []string
}

// NewIssueOptions creates a default Pull Option
func NewIssueOptions() *IssueOptions {
	return &IssueOptions{
		RangeTime: newRangeTime(),
		State:     "open",
		Sort:      "updated",
		Limit:     20,
	}
}

func (opts *IssueOptions) filterIssue(issue *github.Issue) bool {
	if issue.IsPullRequest() {
		return false
	}

	at := issue.GetUpdatedAt()
	if opts.End.Before(at) {
		return false
	}

	return true
}

func (opts *IssueOptions) beforeStart(issue *github.Issue) bool {
	at := issue.GetUpdatedAt()
	return opts.Start.After(at)
}

func (c *Client) GetIssue(ctx context.Context, owner string, repo string, id int) (*github.Issue, error) {
	r, _, err := c.c.Issues.Get(ctx, owner, repo, id)
	return r, err
}

func (c *Client) ListIssueComments(ctx context.Context, owner string, repo string, number int) ([]*github.IssueComment, error) {
	var allComments []*github.IssueComment

	opts := github.IssueListCommentsOptions{
		Sort:      "updated",
		Direction: "desc",
	}

	for {
		comments, resp, err := c.c.Issues.ListComments(ctx, owner, repo, number, &opts)
		if err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allComments, nil
}

func (c *Client) ListIssues(ctx context.Context, opts *IssueOptions, repos []Repository) (map[string][]*github.Issue, error) {
	m := make(map[string][]*github.Issue, len(c.cfg.Repos))
	for _, repo := range repos {
		issues, err := c.ListIssuesByRepo(ctx, &repo, opts)
		if err != nil {
			return nil, err
		}
		m[fmt.Sprintf("%s/%s", repo.Owner, repo.Name)] = issues
	}
	return m, nil
}

func (c *Client) ListIssuesByRepo(ctx context.Context, repo *Repository, opts *IssueOptions) ([]*github.Issue, error) {
	listOpts := github.IssueListByRepoOptions{
		State:       opts.State,
		Sort:        "updated",
		Direction:   "desc",
		Since:       opts.Start,
		ListOptions: github.ListOptions{Page: 0, PerPage: opts.Limit},
	}

	var allIssues []*github.Issue
LOOP:
	for {
		issues, resp, err := c.c.Issues.ListByRepo(ctx, repo.Owner, repo.Name, &listOpts)
		if err != nil {
			return nil, err
		}

		for _, issue := range issues {
			if !opts.filterIssue(issue) {
				continue
			}

			if !filterUsers(issue.Assignees, opts.Assignees) {
				continue
			}

			if opts.beforeStart(issue) {
				break LOOP
			}

			allIssues = append(allIssues, issue)
		}

		if len(allIssues) >= opts.Limit {
			allIssues = allIssues[0:opts.Limit]
			break
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}
