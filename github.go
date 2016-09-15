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
	User      string
	State     string
	New       bool
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
		State:       filter.State,
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

	if filter.User != "" {
		f := filtered
		filtered = []*github.PullRequest{}
		for _, pull := range f {
			if pull.User != nil && *(pull.User.Login) == filter.User {
				filtered = append(filtered, pull)
			}
		}
	}

	if filter.New {
		f := filtered
		filtered = []*github.PullRequest{}
		for _, pull := range f {
			if time.Since(*(pull.CreatedAt)) <= 24*time.Hour {
				filtered = append(filtered, pull)
			}
		}
	}

	return filtered, err
}

// GetIssues -
func (g *GitHubClient) GetIssues(repo Repo, filter Filter) (issues []*github.Issue, err error) {
	issues, _, err = g.client.Issues.ListByRepo(repo.Org, repo.Name, &github.IssueListByRepoOptions{
		State:       filter.State,
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

	if filter.User != "" {
		f := filtered
		filtered = []*github.Issue{}
		for _, issue := range f {
			if issue.User != nil && *(issue.User.Login) == filter.User {
				filtered = append(filtered, issue)
			}
		}
	}

	if filter.New {
		f := filtered
		filtered = []*github.Issue{}
		for _, issue := range f {
			if time.Since(*(issue.CreatedAt)) <= 24*time.Hour {
				filtered = append(filtered, issue)
			}
		}
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
