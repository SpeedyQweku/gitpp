package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/v58/github"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/oauth2"
)

type Config struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Data struct {
	username string
	token    string
	sortlist string
	visilist string
	affilist string
	filePath string
	repoName goflags.StringSlice
	public   bool
	private  bool
	listrepo bool
}

var usd Data

func init() {
	// Get the user's home directory
	hDir, err := os.UserHomeDir()
	if err != nil {
		gologger.Error().Msgf("Error getting user's home directory: %s", err)
		return
	}
	// Specify the folder path and file name
	folderPath := hDir + "/.config/gitpp"
	fileName := "config.json"
	usd.filePath = folderPath + "/" + fileName
	// Check if the folder exists, and create it if not
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			gologger.Error().Msgf("Error creating folder: %s", err)
			return
		}
	}
	// Check if the file exists
	if _, err := os.Stat(usd.filePath); os.IsNotExist(err) {
		// File does not exist, create a default config
		config := Config{
			Username: "",
			Token:    "",
		}
		// Marshal the config to JSON
		configJSON, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			gologger.Error().Msgf("Error marshaling config: %s", err)
			return
		}
		// Write the config to the file
		err = os.WriteFile(usd.filePath, configJSON, 0644)
		if err != nil {
			gologger.Error().Msgf("Error writing config file: %s", err)
			return
		}
	}
}

func readConfig() (username, token string) {
	// Read the config file
	configFile, err := os.ReadFile(usd.filePath)
	if err != nil {
		gologger.Fatal().Msgf("Error reading config file: %s", err)
		return "", ""
	}

	// Unmarshal the JSON into a Config struct
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		gologger.Fatal().Msgf("Error unmarshaling config: %s", err)
		return "", ""
	}

	return config.Username, config.Token
}

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

func listRepos(token, sortlist, visilist, affilist string) []*github.Repository {
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

func runner(username, token string) {
	if usd.listrepo && len(usd.repoName) == 0 {
		repos := listRepos(token, usd.sortlist, usd.visilist, usd.affilist)
		printTable(repos)
	}

	if len(usd.repoName) > 0 && !usd.listrepo {
		repoNames := usd.repoName
		if usd.private {
			for _, reponame := range repoNames {
				gitPP(token, username, reponame, usd.private)
			}
		} else if usd.public {
			usd.public = false
			for _, reponame := range repoNames {
				gitPP(token, username, reponame, usd.private)
			}
		} else {
			gologger.Fatal().Msg("Please specify Private/Public, -h/--help for help.")
		}
	}

	if usd.listrepo && len(usd.repoName) > 0 {
		repoNames := usd.repoName
		if usd.private {
			for _, reponame := range repoNames {
				gitPP(token, username, reponame, usd.private)
			}
		} else if usd.public {
			usd.public = false
			for _, reponame := range repoNames {
				gitPP(token, username, reponame, usd.private)
			}
		} else {
			gologger.Fatal().Msg("Please specify Private/Public, -h/--help for help.")
		}
		repos := listRepos(token, usd.sortlist, usd.visilist, usd.affilist)
		printTable(repos)
	}
}

func main() {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("gitpp, helps to make your github repo public/private...")
	flagSet.CreateGroup("input", "INPUT",
		flagSet.StringVarP(&usd.username, "username", "u", "", "GitHub username"),
		flagSet.StringVarP(&usd.token, "token", "t", "", "GitHub personal access token"),
		flagSet.StringSliceVarP(&usd.repoName, "repo", "r", nil, "Repository name", goflags.CommaSeparatedStringSliceOptions),
		flagSet.BoolVarP(&usd.public, "public", "pub", false, "Make a repo public"),
		flagSet.BoolVarP(&usd.private, "private", "pvt", false, "Make a repo private"),
	)

	flagSet.CreateGroup("probes", "PROBES",
		flagSet.StringVarP(&usd.sortlist, "sort", "s", "update", "The property to sort the results by. \033[33m[created, updated, pushed, full_name]\033[0m"),
		flagSet.StringVarP(&usd.visilist, "vis", "v", "all", "Limit results to repositories with the specified visibility. \033[33m[all, public, private]\033[0m"),
		flagSet.StringVarP(&usd.affilist, "affil", "a", "owner", "List repos of given affiliation. \033[33m[owner,collaborator,organization_member]\033[0m"),
		flagSet.BoolVarP(&usd.listrepo, "list", "l", false, "List all your repos"),
	)
	_ = flagSet.Parse()

	if usd.username != "" && usd.token != "" {
		gologger.Print().Msgf("[\033[33mWRN\033[0m] Kindly Use The Config File [%s]", usd.filePath)
		runner(usd.username, usd.token)
	} else {
		configUsername, configToken := readConfig()
		if configUsername != "" && configToken != "" {
			runner(configUsername, configToken)
		} else {
			gologger.Fatal().Msgf("Go Edit [%s] Or -h/--help For Help.", usd.filePath)
		}
	}
}
