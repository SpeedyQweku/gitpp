package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/v58/github"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/oauth2"
)

func gitPP(token, username, repoName string, chPP bool) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, _, err := client.Repositories.Edit(ctx, username, repoName, &github.Repository{
		Private: github.Bool(chPP),
	})
	if err != nil {
		gologger.Fatal().Msgf("Error updating repository visibility: %v", err)
	}

	visibility := map[bool]string{true: "Private", false: "Public"}[chPP]
	gologger.Info().Msgf("Repository successfully set to \033[95m%s.\033[0m", visibility)
	gologger.Info().Msgf("Repository URL: https://github.com/%s/%s", username, repoName)
}

func listRepos(token, username, sortlist, visilist, affilist string) []*github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var allRepos []*github.Repository

	opt := &github.RepositoryListByAuthenticatedUserOptions{Affiliation: affilist, Visibility: visilist, Sort: sortlist}
	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
	if err != nil {
		log.Fatalf("Error listing repositories: %v", err)
	}
	allRepos = append(allRepos, repos...)
	return allRepos
}

func printTable(repos []*github.Repository) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	gologger.Info().Msgf("Your List Of Repository And Visibility\n")
	fmt.Fprintln(w, "----------\t ----------")
	fmt.Fprintln(w, "Repository\t Visibility")
	fmt.Fprintln(w, "----------\t ----------")
	privateCount, publicCount := 0, 0
	for _, repo := range repos {
		visibility := map[bool]string{true: "\033[31mPrivate\033[0m", false: "\033[34mPublic\033[0m"}[*repo.Private]
		fmt.Fprintf(w, "%s\t %s\n", repo.GetFullName(), visibility)

		if *repo.Private {
			privateCount++
		} else {
			publicCount++
		}
	}
	fmt.Fprintf(w, "\033[36m\nTotal Repo:\033[0m %d \033[36mPrivate Repo:\033[0m %d \033[36mPublic Repo:\033[0m %d\n", len(repos), privateCount, publicCount)
	_ = w.Flush()
}

func main() {
	var (
		username string
		token    string
		repoName string
		sortlist string
		visilist string
		affilist string
		public   bool
		private  bool
		listrepo bool
	)

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("gitpp, helps to make your git repo public/private...")
	flagSet.CreateGroup("input", "INPUT",
		flagSet.StringVarP(&username, "username", "u", "", "GitHub username"),
		flagSet.StringVarP(&token, "token", "t", "", "GitHub personal access token"),
		flagSet.StringVarP(&repoName, "repo", "r", "", "Repository name"),
		flagSet.BoolVarP(&public, "public", "pub", false, "Makes repo public"),
		flagSet.BoolVarP(&private, "private", "pvt", false, "Makes repo private"),
	)

	flagSet.CreateGroup("probes", "PROBES",
		flagSet.StringVarP(&sortlist, "sort", "s", "update", "The property to sort the results by. \033[33m[created, updated, pushed, full_name]\033[0m"),
		flagSet.StringVarP(&visilist, "vis", "v", "all", "Limit results to repositories with the specified visibility. \033[33m[all, public, private]\033[0m"),
		flagSet.StringVarP(&affilist, "affil", "a", "owner", "List repos of given affiliation. \033[33m[owner,collaborator,organization_member]\033[0m"),
		flagSet.BoolVarP(&listrepo, "list", "l", false, "List all your repo"),
	)
	_ = flagSet.Parse()

	if username == "" && token == "" {
		gologger.Fatal().Msg("Error Please provide all required arguments,-h/--help for help.")
	}

	if listrepo && repoName == "" {
		repos := listRepos(token, username, sortlist, visilist, affilist)
		printTable(repos)
	}

	if repoName != "" && !listrepo {
		if private {
			gitPP(token, username, repoName, private)
		} else if public {
			public = false
			gitPP(token, username, repoName, public)
		} else {
			gologger.Fatal().Msg("Please specify Private/Public, -h/--help for help.")
		}
	}

	if listrepo && repoName != "" {
		if private {
			gitPP(token, username, repoName, private)
		} else if public {
			public = false
			gitPP(token, username, repoName, public)
		} else {
			gologger.Fatal().Msg("Please specify Private/Public, -h/--help for help.")
		}
		repos := listRepos(token, username, sortlist, visilist, affilist)
		printTable(repos)
	}
}
