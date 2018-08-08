package main

const TimeFormat string = "2006-01-02 15:04:05"

func filterRepo(allRepos []Repository, args []string) []Repository {
	if len(args) == 2 {
		// use specail owner and repo
		return []Repository{
			{Owner: args[0], Name: args[1]},
		}
	} else if len(args) == 1 {
		// only owner, filter repos by owner
		var repos []Repository
		for _, repo := range allRepos {
			if repo.Owner == args[0] {
				repos = append(repos, repo)
			}
		}
		return repos
	}

	return allRepos
}
