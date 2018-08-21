package main

// TimeFormat is the foramt for time output
const TimeFormat string = "2006-01-02 15:04:05"

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
