# gitpp

gitpp, helps to make your github repo public/private

## Installation

```bash
go install github.com/SpeedyQweku/gitpp@v0.0.1
```

## Config

- added a json config at ~/.config/gitpp/config.json

contains username and github token to be used with the binary
a default json config is generated if one doesn't exist

```bash
{
    "username":"YOUR_GITHUB_USERNAME",
    "token" : "YOUR_GITHUB_TOKEN"
}
```

## Usage

```bash
gitpp, helps to make your github repo public/private...

INPUT:
   -u, -username string  GitHub username
   -t, -token string     GitHub personal access token
   -r, -repo string      Repository name
   -pub, -public         Make a repo public
   -pvt, -private        Make a repo private

PROBES:
   -s, -sort string   The property to sort the results by. [created, updated, pushed, full_name] (default "update")
   -v, -vis string    Limit results to repositories with the specified visibility. [all, public, private] (default "all")
   -a, -affil string  List repos of given affiliation. [owner,collaborator,organization_member] (default "owner")
   -l, -list          List all your repos
```
