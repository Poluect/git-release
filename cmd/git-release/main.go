package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/blang/semver"
	"github.com/poluect/git-release/cmd/git-release/config"
	release "github.com/poluect/git-release/internal"
)

// ResOutput shows output of release operation.
type ResOutput struct {
	Err        error
	NewVersion string
}

var cfg = config.GetConfig()

func main() {
	versionFlag := flag.String("version", "", "semantic version of new release. May be one of patch,minor,major, or specific version e.g. 1.3.5")
	repoFlag := flag.String("repo", "", "repository name to build new release")
	flag.Parse()

	var (
		valid bool
		v     = *versionFlag
	)
	switch v {
	case "patch", "minor", "major":
		valid = true
	case "":
	default:
		_, err := semver.Parse(v)
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
		valid = true
	}
	if !valid {
		log.Fatalf("cannot parse version param %s. Should be one of patch,minor,major, or specific version e.g. 1.3.5", v)
	}
	if *repoFlag == "" {
		log.Fatal("repository name cannot be empty")
	}

	resChan := make(chan ResOutput)
	go releaseNewVersion(context.Background(), resChan, *repoFlag, v)
	out := <-resChan
	if out.Err != nil {
		log.Fatalf("release failed: %v\n", out.Err)
	}

	log.Printf("new release version: %s\n", out.NewVersion)
}

func releaseNewVersion(ctx context.Context, resChan chan ResOutput, repo, version string) {
	release.InitGithubClient(ctx, cfg.GithubToken)

	log.Printf("specified release version: %s\n", version)

	_, err := release.GetRepoInfo(ctx, cfg.OrganizationName, repo)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail get repo (%s): %v", repo, err)}
		return
	}

	latestTag, err := release.GetLatestVersion(ctx, cfg.OrganizationName, repo)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail get latest version: %v", err)}
		return
	}
	log.Printf("latest tag found: %s\n", latestTag)

	newVersion, err := release.CreateNewVersion(version, latestTag)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail create new version: %v", err)}
		return
	}
	log.Printf("new version is: %s\n", newVersion)

	resChan <- ResOutput{NewVersion: newVersion}
}
