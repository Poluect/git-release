package internal

import (
	"context"
	"fmt"

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

// GetBranch gets branch of a specific repo.
func GetBranch(ctx context.Context, org, repoName, branchName string) (*github.Branch, error) {
	branch, _, err := client.Repositories.GetBranch(ctx, org, repoName, branchName)
	return branch, err
}

// CreateBranch creates new branch in a github.
func CreateBranch(ctx context.Context, org, repoName, branchName string, sha *string) (*github.Reference, error) {
	refName := fmt.Sprintf("refs/heads/%s", branchName)

	ref, _, err := client.Git.CreateRef(ctx, org, repoName, &github.Reference{
		Object: &github.GitObject{
			SHA: sha,
		},
		Ref: &refName,
	})

	return ref, err
}

// DeleteBranch deletes branch.
func DeleteBranch(ctx context.Context, org, repoName, sha string) error {
	_, err := client.Git.DeleteRef(ctx, org, repoName, sha)
	return err
}

// CreateTag creates tag for a specific revision.
func CreateTag(ctx context.Context, org, repoName, tagName string, sha *string) (*github.Tag, error) {
	objType := "commit"
	tag, _, err := client.Git.CreateTag(ctx, org, repoName, &github.Tag{
		Message: &tagName,
		Tag:     &tagName,
		Object: &github.GitObject{
			SHA:  sha,
			Type: &objType,
		},
	})

	return tag, err
}

// CreatePR creates new pull request.
func CreatePR(ctx context.Context, org, repoName, branchFrom, branchTo, title string) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Create(ctx, org, repoName, &github.NewPullRequest{
		Base:  &branchTo,
		Head:  &branchFrom,
		Title: &title,
	})

	return pr, err
}
