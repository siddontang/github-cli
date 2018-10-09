package main

import (
	"context"

	"github.com/google/go-github/github"
)

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
