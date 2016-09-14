package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubClient -
type GitHubClient struct {
	client *github.Client
}

// Repo -
type Repo struct {
	Org  string
	Name string
}

// Filter -
type Filter struct {
	Milestone string
}

func (r Repo) String() string {
	return r.Org + "/" + r.Name
}

// NewGitHubClient -
func NewGitHubClient() *GitHubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	return &GitHubClient{client: client}
}

// FindTeamID -
func (g *GitHubClient) FindTeamID(org string, teamName string) (ID int, err error) {
	var last bool
	var next int
	var teams []*github.Team

	for true {
		t, resp, err := g.client.Organizations.ListTeams(org, &github.ListOptions{Page: next})
		if err != nil {
			return -1, err
		}
		teams = append(teams, t...)
		next = resp.NextPage

		// the _next_ iteration is going to be the last one...
		if last {
			break
		}
		last = resp.NextPage == resp.LastPage
	}

	for _, team := range teams {
		if *(team.Name) == teamName || *(team.Slug) == teamName {
			return *(team.ID), nil
		}
	}
	return -1, fmt.Errorf("Team '%s' not found in org '%s'", teamName, org)
}

// GetRepos -
func (g *GitHubClient) GetRepos(org string, teamName string) (repoList []Repo, err error) {
	teamID, err := g.FindTeamID(org, teamName)
	if err != nil {
		return nil, err
	}

	repos, _, err := g.client.Organizations.ListTeamRepos(teamID, &github.ListOptions{PerPage: 200})
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if *(repo.Owner.Login) == org {
			repoList = append(repoList, Repo{*(repo.Owner.Login), *(repo.Name)})
		}
	}
	return repoList, err
}

// GetPRs -
func (g *GitHubClient) GetPRs(repo Repo, filter Filter) (pulls []*github.PullRequest, err error) {
	pulls, _, err = g.client.PullRequests.List(repo.Org, repo.Name, &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 200},
	})
	var filtered []*github.PullRequest
	if filter.Milestone != "" {
		for _, pull := range pulls {
			if pull.Milestone != nil && *(pull.Milestone.Title) == filter.Milestone {
				filtered = append(filtered, pull)
			}
		}
	} else {
		filtered = pulls
	}
	return filtered, err
}

// GetIssues -
func (g *GitHubClient) GetIssues(repo Repo, filter Filter) (issues []*github.Issue, err error) {
	issues, _, err = g.client.Issues.ListByRepo(repo.Org, repo.Name, &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 200},
	})
	var filtered []*github.Issue
	if filter.Milestone != "" {
		for _, issue := range issues {
			if issue.Milestone != nil && *(issue.Milestone.Title) == filter.Milestone {
				filtered = append(filtered, issue)
			}
		}
	} else {
		filtered = issues
	}
	return filtered, err
}

// DisplayPullRequests -
func (g *GitHubClient) DisplayPullRequests(w io.Writer, repo Repo, pulls []*github.PullRequest) {
	for _, p := range pulls {
		var number = *(p.Number)

		var updated = HumanDuration(time.Since(*(p.UpdatedAt)))
		var contributor = *(p.User.Login)

		var milestone string
		if p.Milestone != nil {
			if time.Since(*(p.Milestone.DueOn)).Seconds() > 0 {
				milestone = rgbterm.String(*(p.Milestone.Title), 255, 0, 0, 0, 0, 0)
			} else {
				milestone = rgbterm.String(*(p.Milestone.Title), 192, 192, 192, 0, 0, 0)
			}
		} else {
			milestone = rgbterm.String("", 192, 192, 192, 0, 0, 0)
		}

		var title = *(p.Title)

		pr := fmt.Sprintf("%s/%s#%d", repo.Org, repo.Name, number)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nc(pr), nc(updated), nc(contributor), milestone, nc(title))
	}
}

// DisplayIssues -
func (g *GitHubClient) DisplayIssues(w io.Writer, repo Repo, issues []*github.Issue) {
	for _, issue := range issues {
		var number = *(issue.Number)

		var updated = HumanDuration(time.Since(*(issue.UpdatedAt)))
		var contributor = *(issue.User.Login)

		var milestone string
		if issue.Milestone != nil {
			if time.Since(*(issue.Milestone.DueOn)).Seconds() > 0 {
				milestone = rgbterm.String(*(issue.Milestone.Title), 255, 0, 0, 0, 0, 0)
			} else {
				milestone = rgbterm.String(*(issue.Milestone.Title), 192, 192, 192, 0, 0, 0)
			}
		} else {
			milestone = rgbterm.String("", 192, 192, 192, 0, 0, 0)
		}

		var title = *(issue.Title)

		i := fmt.Sprintf("%s/%s#%d", repo.Org, repo.Name, number)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nc(i), nc(updated), nc(contributor), milestone, nc(title))
	}
}

// shortcut for "no" colour
func nc(s string) string {
	return rgbterm.String(s, 255, 255, 255, 0, 0, 0)
}

// HumanDuration returns a human-readable approximation of a duration
// This function is taken from the Docker project, and slightly modified
// to cap units at days.
// (eg. "About a minute", "4 hours ago", etc.)
// (c) 2013 Docker, inc. and the Docker authors (http://docker.io)
func HumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours()); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}