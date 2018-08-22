# Github CLI

A simple tool to help me manages too many Github repositories.

## Prepare

1. Create a `.github-cli` directory in your home directory. 

    ```bash
    mkdir -p ~/.github-cli
    ```

2. [Option] Create a personal [token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) and save this token to `~/.github-cli/token`.

3. Create a configuration in `~/.github-cli/config.toml`, see [example](./config.toml) here.

## Usage

### Trending

```bash
github-cli trending
# Trending go language
github-cli trending go
# Trending go language in this week
github-cli trending go --time weekly
```

### Pull

```bash
github-cli pulls
# List pull requests of repositories configured in tikv organization
github-cli pulls --owner tikv
# List pull requests of tikv repository in tikv organization 
github-cli pulls tikv
# or use github-cli pulls tikv/tikv
# See one pull request
github-cli pull tikv 3344
```

### Issue

```bash
github-cli issues
# List issues of repositories configured in tikv organization 
github-cli issues --owner tikv
# List issues of tikv repository in tikv organization 
github-cli issues tikv
# or use github-cli issues tikv/tikv
# See one issue 
github-cli issue tikv 3355
```

### Event

```bash
# List events for the account in the configuration file
github-cli events 
# List events for a specified user
github-cli events siddontang
```