package internal

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
)

// GetLatestVersion gets latest tag
func GetLatestVersion(ctx context.Context, org, repo string) (latestTag string, err error) {
	refs, response, err := client.Git.ListRefs(ctx, org, repo, &github.ReferenceListOptions{
		Type: "tags",
	})
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			err = nil
		}
		return
	}
	if len(refs) == 0 {
		return
	}

	for i := len(refs) - 1; i >= 0; i-- {
		tag := filepath.Base(*refs[len(refs)-1].Ref)
		if !strings.HasPrefix(tag, "v") || len(tag) <= len("v") {
			continue
		}

		latestTag = tag[len("v"):]
		break
	}

	return
}

// CreateNewVersion creates new semantic version.
func CreateNewVersion(version, prevVersion string) (res string, err error) {
	var (
		pv, v semver.Version
	)

	if prevVersion != "" {
		pv, err = semver.Make(prevVersion)
		if err != nil {
			err = fmt.Errorf("cannot parse previous version %s: %v", prevVersion, err)
			return
		}
		v = pv
	}

	switch version {
	case "patch":
		v.Patch++
		res = v.String()
	case "minor":
		v.Patch = 0
		v.Minor++
		res = v.String()
	case "major":
		v.Patch = 0
		v.Minor = 0
		v.Major++
		res = v.String()
	default:
		v, err = semver.Make(version)
		if err != nil {
			err = fmt.Errorf("cannot parse new version %s: %v", version, err)
			return
		}
		if prevVersion != "" && !v.GT(pv) {
			err = fmt.Errorf("new version (%s) should be greater than previous one (%s)", v, pv)
			return
		}

		res = v.String()
	}

	return
}
