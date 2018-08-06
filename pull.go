package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
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
	n := time.Now().UTC()
	return &PullOptions{
		Start: n.Add(-14 * 24 * time.Hour),
		End:   n,
		State: "closed",
		Sort:  "updated",
		Limit: 20,
	}
}

// filterPull checks whether pull meets the options
func (opt *PullOptions) filterPull(pull *github.PullRequest) bool {
	if opt.State != pull.GetState() {
		return false
	}

	return true
}

func (opt *PullOptions) beforeStart(pull *github.PullRequest) bool {
	at := pull.GetUpdatedAt()
	return opt.Start.After(at)
}

func (opt *PullOptions) afterEnd(pull *github.PullRequest) bool {
	at := pull.GetUpdatedAt()
	return opt.End.Before(at)
}

func (c *Client) ListPulls(ctx context.Context, opts *PullOptions) (map[string][]*github.PullRequest, error) {
	m := make(map[string][]*github.PullRequest, len(c.cfg.Repos))
	for _, repo := range c.cfg.Repos {
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

			if opts.afterEnd(pull) {
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

const TimeFormat string = "2006-01-02 15:04:05"

var (
	pullState     string
	pullLimit     int
	pullEndTime   string
	pullOffsetDur string
)

func newPullCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "pull",
		Short: "Github CLI for pull",
		Args:  cobra.MinimumNArgs(0),
		Run:   runPullCommandFunc,
	}
	m.Flags().StringVar(&pullState, "state", "open", "PR state: open or closed")
	m.Flags().IntVar(&pullLimit, "limit", 100, "Maximum pull limit for a repository")
	m.Flags().StringVar(&pullEndTime, "end", "", fmt.Sprintf("Pull End Time, format is %s", TimeFormat))
	m.Flags().StringVar(&pullOffsetDur, "offset", "-336h", "Pull offset, if > 0, the time range is [start, start + offset], if < 0, the time range is [end - offset, end]")
	return m
}

func runPullCommandFunc(cmd *cobra.Command, args []string) {
	opts := NewPullOptions()
	opts.State = pullState
	opts.Limit = pullLimit

	if len(pullEndTime) > 0 {
		end, err := time.Parse(TimeFormat, pullEndTime)
		perror(err)
		opts.End = end
	}

	if len(pullOffsetDur) > 0 {
		d, err := time.ParseDuration(pullOffsetDur)
		perror(err)
		if d > 0 {
			opts.End = opts.Start.Add(d)
		} else if d < 0 {
			opts.Start = opts.End.Add(d)
		}
	}

	m, err := globalClient.ListPulls(globalCtx, opts)
	perror(err)

	for repo, pulls := range m {
		fmt.Println(repo)
		for _, pull := range pulls {
			fmt.Printf("%s %s %s\n", pull.GetUpdatedAt().Format(TimeFormat), pull.GetHTMLURL(), pull.GetTitle())
		}
	}
}
