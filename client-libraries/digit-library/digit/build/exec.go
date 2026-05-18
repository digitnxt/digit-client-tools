package build

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CmdResult struct {
	Stdout string
	Stderr string
}

func RunCommand(ctx context.Context, logger *slog.Logger, timeout time.Duration, name string, args ...string) (*CmdResult, error) {
	return RunCommandInDir(ctx, logger, timeout, "", name, args...)
}

func RunCommandInDir(ctx context.Context, logger *slog.Logger, timeout time.Duration, dir string, name string, args ...string) (*CmdResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Info("exec", "cmd", strings.Join(append([]string{name}, args...), " "))

	if err := cmd.Run(); err != nil {
		return &CmdResult{Stdout: stdout.String(), Stderr: stderr.String()}, fmt.Errorf("command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return &CmdResult{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}

func RunCommandWithEnvInDir(ctx context.Context, logger *slog.Logger, timeout time.Duration, dir string, env []string, name string, args ...string) (*CmdResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Info("exec", "cmd", strings.Join(append([]string{name}, args...), " "))

	if err := cmd.Run(); err != nil {
		return &CmdResult{Stdout: stdout.String(), Stderr: stderr.String()}, fmt.Errorf("command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return &CmdResult{Stdout: stdout.String(), Stderr: stderr.String()}, nil
}
