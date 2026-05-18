package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func BuildImage(ctx context.Context, logger *slog.Logger, repoPath, dockerfilePath, repoName, sha, workDir, orgOverride string) (string, error) {
	org := orgOverride
	if org == "" {
		org = os.Getenv("DOCKER_ORG")
	}
	if org == "" {
		return "", errors.New("DOCKER_ORG is required")
	}

	image := fmt.Sprintf("%s/%s:%s", org, repoName, sha)
	if workDir == "" {
		workDir = "."
	}
	args := []string{"build", "--build-arg", fmt.Sprintf("WORK_DIR=%s", workDir), "-f", filepath.Base(dockerfilePath), "-t", image, "."}
	_, err := RunCommandWithEnvInDir(ctx, logger, 20*time.Minute, repoPath, []string{"DOCKER_BUILDKIT=1"}, "docker", args...)
	if err != nil {
		return "", err
	}
	return image, nil
}
