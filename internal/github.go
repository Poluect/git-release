package internal

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var client *github.Client

// InitGithubClient inits github client.
func InitGithubClient(ctx context.Context, token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}

// GetRepoInfo gets repo info.
func GetRepoInfo(ctx context.Context, org, repoName string) (*github.Repository, error) {
	repo, _, err := client.Repositories.Get(ctx, org, repoName)
	return repo, err
}

// CreateBranch creates new branch in a github.
func CreateBranch(branchName string) error {
	return nil
}

// CreatePR creates new pull request.
func CreatePR(from, to string) error {
	return nil
}
