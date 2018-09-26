package main

import (
	"strings"
	"time"

	"github.com/google/go-github/github"
)

// TimeFormat is the foramt for time output
const TimeFormat string = "2006-01-02 15:04:05"

// RangeTime is a time range in [start, end]
type RangeTime struct {
	Start time.Time
	End   time.Time
}

func newRangeTime() RangeTime {
	n := time.Now()
	return RangeTime{
		Start: n.Add(-7 * 24 * time.Hour),
		End:   n,
	}
}

func (r *RangeTime) adjust(sinceTime string, offsetDur string) {
	if len(sinceTime) > 0 {
		end, err := time.Parse(TimeFormat, sinceTime)
		perror(err)
		r.End = end
	}

	d, err := time.ParseDuration(offsetDur)
	perror(err)
	r.Start = r.End.Add(d)
	if r.Start.After(r.End) {
		r.Start, r.End = r.End, r.Start
	}
}

func adjustRepoName(owner string, args []string) (string, string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	if len(name) == 0 {
		return owner, name
	}

	for i := 0; i < len(name); i++ {
		if name[i] == '/' {
			// The name has already been the format of owner/name
			return name[0:i], name[i+1:]
		}
	}

	return owner, name
}

func findRepo(c *Config, args []string) Repository {
	owner, name := adjustRepoName("", args)
	if repo := c.FindRepo(owner, name); repo != nil {
		return *repo
	}

	// use specail owner and repo
	return Repository{Owner: owner, Name: name}

}

func filterRepo(c *Config, owner string, args []string) []Repository {
	var name string
	owner, name = adjustRepoName(owner, args)
	if len(name) == 0 && len(owner) == 0 {
		return c.Repos
	} else if len(name) == 0 {
		// only owner, filter repos by owner
		var repos []Repository
		for _, repo := range c.Repos {
			if repo.Owner == owner {
				repos = append(repos, repo)
			}
		}
		return repos
	}

	// only name
	if r := c.FindRepo(owner, name); r != nil {
		return []Repository{*r}
	}

	// use specail owner and repo
	return []Repository{
		{Owner: owner, Name: name},
	}
}

func splitUsers(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	return strings.Split(s, ",")
}

func filterUsers(users []*github.User, names []string) bool {
	if len(names) == 0 {
		return true
	}

	for _, name := range names {
		for _, user := range users {
			if user.GetLogin() == name {
				return true
			}
		}
	}

	return false
}
