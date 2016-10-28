package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hairyhenderson/teams/version"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("teams", "Manage GitHub issues and pull requests for a team's repositories")
	app.Version("v version", version.Version)

	teamOpt := app.StringOpt("t team", "", "Team name to manage (in org/teamname format)")

	g := NewGitHubClient()
	var repoList []Repo

	// app.Before =
	app.Before = func() {
		t := strings.SplitN(*teamOpt, "/", 2)
		org := t[0]
		team := t[1]

		repoList, _ = g.GetRepos(org, team)
	}

	app.Command("pulls", "Displays open PRs for the given repo(s)", func(cmd *cli.Cmd) {
		stateArg := cmd.StringOpt("state", "open", "Display PRs based on their state")
		milestoneArg := cmd.StringOpt("m milestone", "", "Milestone to filter on")
		userArg := cmd.StringOpt("u user", "", "Display only PRs from <user>")
		newArg := cmd.BoolOpt("new", false, "Display PRs opened within the past 24 hours")

		cmd.Action = func() {
			w := tabwriter.NewWriter(os.Stdout, 4, 2, 2, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nc("PULL REQUEST"), nc("LAST UPDATED"), nc("CONTRIBUTOR"), nc("MILESTONE"), nc("TITLE"))
			filter := Filter{
				State: *(stateArg),
				New:   *(newArg),
			}
			if milestoneArg != nil {
				filter.Milestone = *(milestoneArg)
			}
			if userArg != nil {
				filter.User = *(userArg)
			}
			for _, repo := range repoList {
				pulls, _ := g.GetPRs(repo, filter)
				g.DisplayPullRequests(w, repo, pulls)
			}
			w.Flush()
		}
	})

	app.Command("issues", "Displays open Issues for the given repo(s)", func(cmd *cli.Cmd) {
		stateArg := cmd.StringOpt("state", "open", "Display Issues based on their state")
		milestoneArg := cmd.StringOpt("m milestone", "", "Milestone to filter on")
		userArg := cmd.StringOpt("u user", "", "Display only Issues from <user>")
		newArg := cmd.BoolOpt("new", false, "Display Issues opened within the past 24 hours")
		cmd.Action = func() {
			w := tabwriter.NewWriter(os.Stdout, 4, 2, 2, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nc("ISSUE"), nc("LAST UPDATED"), nc("CONTRIBUTOR"), nc("MILESTONE"), nc("TITLE"))
			filter := Filter{
				State: *(stateArg),
				New:   *(newArg),
			}
			if milestoneArg != nil {
				filter.Milestone = *(milestoneArg)
			}
			if userArg != nil {
				filter.User = *(userArg)
			}
			for _, repo := range repoList {
				issues, _ := g.GetIssues(repo, filter)
				g.DisplayIssues(w, repo, issues)
			}
			w.Flush()
		}
	})

	app.Command("list", "Lists all repos owned by the given team", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			w := tabwriter.NewWriter(os.Stdout, 4, 2, 2, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\n", nc("REPO"), nc("DESCRIPTION"))
			g.DisplayRepos(w, repoList)
			w.Flush()
		}
	})

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
