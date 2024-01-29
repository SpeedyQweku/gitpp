package main

import (
	"context"

	"github.com/google/go-github/v35/github"
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
	gologger.Info().Msgf("Repository successfully set to \033[31m%s.\033[0m", visibility)
	gologger.Info().Msgf("Repository URL: https://github.com/%s/%s", username, repoName)
}

func main() {
	var username string
	var token string
	var repoName string
	var public bool
	var private bool

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("gitpp, helps to make your git repo public/private.")
	flagSet.CreateGroup("input", "INPUT",
		flagSet.StringVarP(&username, "username", "u", "", "GitHub username"),
		flagSet.StringVarP(&token, "token", "t", "", "GitHub personal access token"),
		flagSet.StringVarP(&repoName, "repo", "r", "", "Repository name"),
		flagSet.BoolVarP(&public, "public", "pub", false, "Makes repo public"),
		flagSet.BoolVarP(&private, "private", "pvt", false, "Makes repo private"),
	)
	_ = flagSet.Parse()

	if username == "" || token == "" || repoName == "" {
		gologger.Fatal().Msg("Error Please provide all required arguments,-h/--help for help.")
	}

	if private {
		gitPP(token, username, repoName, private)
	} else if public {
		public = false
		gitPP(token, username, repoName, public)
	}else {
		gologger.Fatal().Msg("Please specify Private/Public, -h/--help for help.")
	}
}
