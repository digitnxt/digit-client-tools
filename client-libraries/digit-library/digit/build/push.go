package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

func PushImage(ctx context.Context, logger *slog.Logger, image, userOverride, tokenOverride string) error {
	user := userOverride
	token := tokenOverride
	if user == "" {
		user = os.Getenv("DOCKER_USERNAME")
	}
	if token == "" {
		token = os.Getenv("DOCKER_TOKEN")
	}
	if user == "" || token == "" {
		return errors.New("DOCKER_USERNAME and DOCKER_TOKEN are required")
	}

	if err := dockerLogin(ctx, logger, user, token); err != nil {
		return err
	}

	_, err := RunCommand(ctx, logger, 10*time.Minute, "docker", "push", image)
	return err
}

func dockerLogin(ctx context.Context, logger *slog.Logger, user, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "login", "-u", user, "--password-stdin")
	cmd.Stdin = strings.NewReader(token)

	logger.Info("exec", "cmd", "docker login -u *** --password-stdin")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker login failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
