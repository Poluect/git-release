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
	dryFlag := flag.Bool("dry", false, "specify --dry=true if you want to see what semantic version is going to be released")
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
		log.Fatalf("cannot parse version param (%s). Should be one of patch,minor,major, or specific version e.g. 1.3.5", v)
	}
	if *repoFlag == "" {
		log.Fatal("repository name cannot be empty")
	}

	ctx := context.Background()
	if cfg.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	resChan := make(chan ResOutput)
	go releaseNewVersion(ctx, resChan, *repoFlag, v, *dryFlag)
	out := <-resChan
	if out.Err != nil {
		log.Fatalf("release failed: %v\n", out.Err)
	}

	if !*dryFlag {
		log.Printf("new release version: %s\n", out.NewVersion)
	}
}

func releaseNewVersion(ctx context.Context, resChan chan ResOutput, repo, version string, dry bool) {
	release.InitGithubClient(ctx, cfg.GithubToken)

	log.Printf("specified release version: %s\n", version)

	_, err := release.GetRepoInfo(ctx, cfg.OrganizationName, repo)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail get repo (%s): %v", repo, err)}
		return
	}
	fromBranch, err := release.GetBranch(ctx, cfg.OrganizationName, repo, cfg.BranchReleaseFrom)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail get branch (%s): %v", cfg.BranchReleaseFrom, err)}
		return
	}
	_, err = release.GetBranch(ctx, cfg.OrganizationName, repo, cfg.BranchReleaseTo)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail get branch (%s): %v", cfg.BranchReleaseTo, err)}
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

	releaseBranchName := fmt.Sprintf("release-v%s", newVersion)
	if dry { // just describe what is going to be changed
		log.Printf("branch (%s) is going to be created\n", releaseBranchName)
		log.Printf("PR to branch (%s) from head (%s) is going to be created\n", cfg.BranchReleaseTo, releaseBranchName)
		resChan <- ResOutput{NewVersion: releaseBranchName}
	}

	_, err = release.CreateBranch(ctx, cfg.OrganizationName, repo, releaseBranchName, fromBranch.Commit.SHA)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail create new version: %v", err)}
		return
	}
	log.Printf("branch (%s) has been created\n", releaseBranchName)

	prTitle := fmt.Sprintf("Release v%s", newVersion)
	pr, err := release.CreatePR(ctx, cfg.OrganizationName, repo, releaseBranchName, cfg.BranchReleaseTo, prTitle)
	if err != nil {
		resChan <- ResOutput{Err: fmt.Errorf("fail create pull request: (%v), please create it manually", err)}
		return
	}
	log.Printf("new PR has been created with a state (%s)\n", *pr.State)

	resChan <- ResOutput{NewVersion: releaseBranchName}
}
