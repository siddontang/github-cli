package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// EventOptions is the options for listing events
type EventOptions struct {
	RangeTime

	// maximum number
	Limit int
}

// NewEventOptions creates a default Event Option
func NewEventOptions() *EventOptions {
	return &EventOptions{
		RangeTime: newRangeTime(),
		Limit:     20,
	}
}

func (opts *EventOptions) filterEvent(event *github.Event) bool {
	at := event.GetCreatedAt()
	if opts.End.Before(at) {
		return false
	}

	return true
}

func (opts *EventOptions) beforeStart(event *github.Event) bool {
	at := event.GetCreatedAt()
	return opts.Start.After(at)
}

func (c *Client) ListEventsByUser(ctx context.Context, user string, opts *EventOptions) ([]*github.Event, error) {
	listOpts := github.ListOptions{
		Page: 0, PerPage: opts.Limit,
	}

	var allEvents []*github.Event
LOOP:
	for {
		events, resp, err := c.c.Activity.ListEventsPerformedByUser(ctx, user, true, &listOpts)
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			if !opts.filterEvent(event) {
				continue
			}

			if opts.beforeStart(event) {
				break LOOP
			}

			allEvents = append(allEvents, event)
		}

		if len(allEvents) >= opts.Limit {
			allEvents = allEvents[0:opts.Limit]
			break
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return allEvents, nil
}

func formatEvent(event *github.Event) string {
	// We only care some events now
	payload, err := event.ParsePayload()
	perror(err)

	switch e := payload.(type) {
	case *github.IssuesEvent:
		return fmt.Sprintf("Issues: %s %s %s", e.GetAction(), e.GetIssue().GetHTMLURL(), e.GetIssue().GetTitle())
	case *github.IssueCommentEvent:
		return fmt.Sprintf("IssueComment: %s %s", e.GetIssue().GetHTMLURL(), e.GetComment().GetBody())
	case *github.PullRequestEvent:
		return fmt.Sprintf("Pull: %s %s %s", e.GetAction(), e.GetPullRequest().GetHTMLURL(), e.GetPullRequest().GetTitle())
	case *github.PullRequestReviewCommentEvent:
		return fmt.Sprintf("PullComment: %s %s", e.GetPullRequest().GetHTMLURL(), e.GetComment().GetBody())
	}

	return event.GetType()
}
