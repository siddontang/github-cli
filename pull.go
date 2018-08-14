package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// PullOptions is the options for listing pulls
type PullOptions struct {
	Start time.Time
	End   time.Time
	// open or closed
	State string
	// created, updated, popularity, long-running
	// Default: updated
	Sort string
	// maximum number
	Limit int
}

// NewPullOptions creates a default Pull Option
func NewPullOptions() *PullOptions {
	n := time.Now()
	return &PullOptions{
		Start: n.Add(-14 * 24 * time.Hour),
		End:   n,
		State: "closed",
		Sort:  "updated",
		Limit: 20,
	}
}

// filterPull checks whether pull meets the options
func (opts *PullOptions) filterPull(pull *github.PullRequest) bool {
	at := pull.GetUpdatedAt()
	if opts.End.Before(at) {
		return false
	}

	return true
}

func (opts *PullOptions) beforeStart(pull *github.PullRequest) bool {
	at := pull.GetUpdatedAt()
	return opts.Start.After(at)
}

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

func (c *Client) ListPulls(ctx context.Context, opts *PullOptions, repos []Repository) (map[string][]*github.PullRequest, error) {
	m := make(map[string][]*github.PullRequest, len(c.cfg.Repos))
	for _, repo := range repos {
		pulls, err := c.ListPullsByRepo(ctx, &repo, opts)
		if err != nil {
			return nil, err
		}
		m[fmt.Sprintf("%s/%s", repo.Owner, repo.Name)] = pulls
	}
	return m, nil
}

// ListPullsByRepo lists the pulls by repository
func (c *Client) ListPullsByRepo(ctx context.Context, repo *Repository, opts *PullOptions) ([]*github.PullRequest, error) {
	listOpts := github.PullRequestListOptions{
		State:       opts.State,
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{Page: 0, PerPage: opts.Limit},
	}

	var allPulls []*github.PullRequest
LOOP:
	for {
		pulls, resp, err := c.c.PullRequests.List(ctx, repo.Owner, repo.Name, &listOpts)
		if err != nil {
			return nil, err
		}

		for _, pull := range pulls {
			if !opts.filterPull(pull) {
				continue
			}

			if opts.beforeStart(pull) {
				break LOOP
			}

			allPulls = append(allPulls, pull)
		}

		if len(allPulls) >= opts.Limit {
			allPulls = allPulls[0:opts.Limit]
			break
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.ListOptions.Page = resp.NextPage
	}

	return allPulls, nil
}
