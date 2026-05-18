package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type BuildOptions struct {
	RepoURL  string
	Branch   string
	WorkDir  string
	ImageName      string
	SkipScan       bool
	DockerUsername string
	DockerToken    string
	DockerOrg      string
	Logger   *slog.Logger
	Now      func() time.Time
	NewCtx   func() context.Context
	Stdout   *os.File
}

func RunBuild(opts BuildOptions) error {
	if opts.RepoURL == "" {
		return errors.New("repo url is required")
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	logger := opts.Logger
	if logger == nil {
		logger = newLogger()
	}
	newCtx := opts.NewCtx
	if newCtx == nil {
		newCtx = context.Background
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	start := now()
	ctx := newCtx()

	repoInfo, err := ValidateRepoURL(opts.RepoURL)
	if err != nil {
		return wrapStageError("validate", err)
	}

	workdir, cleanup, err := CloneRepo(ctx, logger, repoInfo, opts.Branch)
	if err != nil {
		return wrapStageError("clone", err)
	}
	defer cleanup()

	logger.Info("repo cloned", "path", workdir)

	var detection Detection
	if opts.WorkDir != "" {
		detection, err = DetectLanguageInDir(workdir, opts.WorkDir)
	} else {
		detection, err = DetectLanguage(workdir)
	}
	if err != nil {
		return wrapStageError("detect", err)
	}

	if opts.WorkDir != "" {
		logger.Info("workdir override applied", "workdir", opts.WorkDir)
	}

	dockerfilePath, err := EnsureDockerfile(ctx, logger, workdir, detection, detection.DockerfilePath)
	if err != nil {
		return wrapStageError("dockerfile", err)
	}

	if err := RunTests(ctx, logger, workdir, detection.Lang, detection.WorkDir); err != nil {
		return wrapStageError("tests", err)
	}

	repoName := strings.ToLower(repoInfo.Name)
	if opts.ImageName != "" {
		repoName = strings.ToLower(opts.ImageName)
	}
	sha, err := GitShortSHA(ctx, logger, workdir)
	if err != nil {
		return wrapStageError("git", err)
	}

	image, err := BuildImage(ctx, logger, workdir, dockerfilePath, repoName, sha, detection.WorkDir, opts.DockerOrg)
	if err != nil {
		return wrapStageError("build", err)
	}

	summary := ScanSummary{Skipped: true}
	if !opts.SkipScan {
		summary, err = ScanImage(ctx, logger, image)
		if err != nil {
			return wrapStageError("scan", err)
		}
	} else {
		logger.Info("scan skipped by flag")
	}

	if err := PushImage(ctx, logger, image, opts.DockerUsername, opts.DockerToken); err != nil {
		return wrapStageError("push", err)
	}

	duration := time.Since(start).Round(time.Second)
	fmt.Fprintf(stdout, "Image pushed: %s\n", image)
	fmt.Fprintln(stdout, summary.String())
	fmt.Fprintf(stdout, "Total time: %s\n", duration)
	return nil
}

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler)
}

func wrapStageError(stage string, err error) error {
	if err == nil {
		return nil
	}

	var se *StageError
	if errors.As(err, &se) {
		return err
	}

	return &StageError{Stage: stage, Err: err}
}
