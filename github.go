package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
)

// SearchOptions is the options for Search Issues
type SearchOptions struct {
	Sort  string
	Order string
	Limit int
}

// SearchIssues provides a common way to search issues or pull requests.
func (c *Client) SearchIssues(ctx context.Context, repos []Repository, opts SearchOptions, queryArgs url.Values) (map[string][]github.Issue, error) {
	m := make(map[string][]github.Issue, len(c.cfg.Repos))
	for _, repo := range repos {
		issues, err := c.SearchIssuesByRepo(ctx, repo, opts, queryArgs)
		if err != nil {
			return nil, err
		}
		m[fmt.Sprintf("%s/%s", repo.Owner, repo.Name)] = issues
	}
	return m, nil
}

// SearchIssuesByRepo provides a common way to search issues or pull requests.
func (c *Client) SearchIssuesByRepo(ctx context.Context, repo Repository, opts SearchOptions, queryArgs url.Values) ([]github.Issue, error) {
	opt := github.SearchOptions{
		Sort:  opts.Sort,
		Order: opts.Order,
	}

	queryArgs.Del("repo")
	queryArgs.Add("repo", repo.String())

	var (
		query bytes.Buffer
		first = true
	)
	for key, values := range queryArgs {
		for _, value := range values {
			if !first {
				query.WriteByte(' ')
			}
			first = false
			query.WriteString(fmt.Sprintf("%s:%s", key, value))
		}
	}

	var allIssues []github.Issue
	for {
		issues, resp, err := c.c.Search.Issues(ctx, query.String(), &opt)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues.Issues...)

		if opts.Limit > 0 && len(allIssues) >= opts.Limit {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}
