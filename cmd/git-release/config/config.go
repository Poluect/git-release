package config

import (
	"os"
)

// Config represents release configuration.
type Config struct {
	GithubToken                        string
	OrganizationName                   string
	BranchReleaseFrom, BranchReleaseTo string
}

var cfg *Config

func init() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = "c96e9c79feb6970dad7c2743318a4a4bdb16cc9f"
	}
	org := os.Getenv("ORGANIZATION_NAME")
	if org == "" {
		org = "cliqueinc"
	}
	releaseFrom := os.Getenv("BRANCH_RELEASE_FROM")
	if releaseFrom == "" {
		releaseFrom = "develop"
	}
	releaseTo := os.Getenv("BRANCH_RELEASE_TO")
	if releaseTo == "" {
		releaseTo = "master"
	}

	cfg = &Config{
		GithubToken:       token,
		OrganizationName:  org,
		BranchReleaseFrom: releaseFrom,
		BranchReleaseTo:   releaseTo,
	}
}

func GetConfig() *Config {
	return cfg
}
