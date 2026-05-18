package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func RunTests(ctx context.Context, logger *slog.Logger, repoPath string, lang Language, workDir string) error {
	testDir := repoPath
	if workDir != "" && workDir != "." {
		testDir = filepath.Join(repoPath, workDir)
	}

	switch lang {
	case LangUnknown:
		logger.Warn("multiple languages detected; skipping tests")
		return nil
	case LangNode:
		return runNodeTests(ctx, logger, testDir)
	case LangGo:
		return runGoTests(ctx, logger, testDir)
	case LangJava:
		return runJavaTests(ctx, logger, testDir)
	default:
		return errors.New("unsupported language for tests")
	}
}

func runNodeTests(ctx context.Context, logger *slog.Logger, repoPath string) error {
	pkgPath := filepath.Join(repoPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return fmt.Errorf("read package.json: %w", err)
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("parse package.json: %w", err)
	}

	script, ok := pkg.Scripts["test"]
	if !ok || strings.TrimSpace(script) == "" {
		logger.Warn("no test script detected; skipping")
		return nil
	}

	_, err = RunCommandInDir(ctx, logger, 10*time.Minute, repoPath, "npm", "test")
	return err
}

func runGoTests(ctx context.Context, logger *slog.Logger, repoPath string) error {
	_, err := RunCommandInDir(ctx, logger, 10*time.Minute, repoPath, "go", "test", "./...")
	return err
}

func runJavaTests(ctx context.Context, logger *slog.Logger, repoPath string) error {
	if fileExists(filepath.Join(repoPath, "pom.xml")) {
		_, err := RunCommandInDir(ctx, logger, 15*time.Minute, repoPath, "mvn", "-f", filepath.Join(repoPath, "pom.xml"), "test")
		return err
	}

	gradleWrapper := filepath.Join(repoPath, "gradlew")
	if fileExists(gradleWrapper) {
		_, err := RunCommandInDir(ctx, logger, 15*time.Minute, repoPath, gradleWrapper, "test")
		return err
	}

	_, err := RunCommandInDir(ctx, logger, 15*time.Minute, repoPath, "gradle", "test")
	return err
}
