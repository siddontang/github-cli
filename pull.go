package main

import (
	"context"

	"github.com/google/go-github/github"
)

func (c *Client) GetPull(ctx context.Context, owner string, repo string, id int) (*github.PullRequest, error) {
	r, _, err := c.c.PullRequests.Get(ctx, owner, repo, id)
	return r, err
}

func (c *Client) ListPullComments(ctx context.Context, owner string, repo string, number int) ([]*github.PullRequestComment, error) {
	opts := github.PullRequestListCommentsOptions{
		Sort:      "updated",
		Direction: "desc",
	}

	var allComments []*github.PullRequestComment
	for {
		comments, resp, err := c.c.PullRequests.ListComments(ctx, owner, repo, number, &opts)
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
