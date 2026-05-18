package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type RepoInfo struct {
	URL  string
	Host string
	Name string
}

var githubSSHPattern = regexp.MustCompile(`^git@github\.com:([A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+)(\.git)?$`)

func ValidateRepoURL(raw string) (*RepoInfo, error) {
	if matches := githubSSHPattern.FindStringSubmatch(raw); len(matches) > 0 {
		name := strings.TrimSuffix(path.Base(matches[1]), ".git")
		return &RepoInfo{URL: raw, Host: "github.com", Name: name}, nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid repo url: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, errors.New("repo url must use http or https")
	}
	if u.Host != "github.com" {
		return nil, errors.New("only github.com is supported")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, errors.New("repo url must be in form https://github.com/org/repo")
	}
	name := strings.TrimSuffix(parts[1], ".git")
	if name == "" {
		return nil, errors.New("repo name missing in url")
	}
	return &RepoInfo{URL: raw, Host: u.Host, Name: name}, nil
}

func CloneRepo(ctx context.Context, logger *slog.Logger, repo *RepoInfo, branch string) (string, func(), error) {
	workdir, err := os.MkdirTemp("", "digit-")
	if err != nil {
		return "", func() {}, fmt.Errorf("create temp dir: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(workdir)
	}

	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repo.URL, workdir)

	_, err = RunCommand(ctx, logger, 2*time.Minute, "git", args...)
	if err != nil {
		cleanup()
		return "", func() {}, err
	}

	return workdir, cleanup, nil
}

func GitShortSHA(ctx context.Context, logger *slog.Logger, repoPath string) (string, error) {
	res, err := RunCommand(ctx, logger, 30*time.Second, "git", "-C", repoPath, "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res.Stdout), nil
}
