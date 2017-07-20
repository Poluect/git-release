package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents release configuration.
type Config struct {
	GithubToken                        string
	OrganizationName                   string
	BranchReleaseFrom, BranchReleaseTo string
	Timeout                            time.Duration
}

var cfg *Config

func init() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("GITHUB token required.")
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
	timeoutSeconds, err := strconv.Atoi(os.Getenv("TIMEOUT_SECONDS"))
	if err != nil || timeoutSeconds < 0 {
		timeoutSeconds = 60
	}

	cfg = &Config{
		GithubToken:       token,
		OrganizationName:  org,
		BranchReleaseFrom: releaseFrom,
		BranchReleaseTo:   releaseTo,
		Timeout:           time.Duration(timeoutSeconds) * time.Second,
	}
}

func GetConfig() *Config {
	return cfg
}
