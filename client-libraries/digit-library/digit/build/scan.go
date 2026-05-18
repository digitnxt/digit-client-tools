package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

type trivyReport struct {
	Results []struct {
		Vulnerabilities []struct {
			Severity string `json:"Severity"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

type ScanSummary struct {
	Critical int
	High     int
	Skipped  bool
}

func (s ScanSummary) String() string {
	if s.Skipped {
		return "Scan summary: SKIPPED"
	}
	return fmt.Sprintf("Scan summary: CRITICAL=%d HIGH=%d", s.Critical, s.High)
}

func ScanImage(ctx context.Context, logger *slog.Logger, image string) (ScanSummary, error) {
	if _, err := exec.LookPath("trivy"); err != nil {
		logger.Warn("trivy not found; skipping scan")
		return ScanSummary{Skipped: true}, nil
	}
	args := []string{"image", "--severity", "CRITICAL,HIGH", "--format", "json", image}
	res, err := RunCommand(ctx, logger, 10*time.Minute, "trivy", args...)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied while trying to connect to the Docker daemon socket") {
			logger.Warn("trivy cannot access docker socket; skipping scan")
			return ScanSummary{Skipped: true}, nil
		}
		return ScanSummary{}, err
	}

	var report trivyReport
	if err := json.Unmarshal([]byte(res.Stdout), &report); err != nil {
		return ScanSummary{}, fmt.Errorf("parse trivy output: %w", err)
	}

	summary := ScanSummary{}
	for _, result := range report.Results {
		for _, vuln := range result.Vulnerabilities {
			switch vuln.Severity {
			case "CRITICAL":
				summary.Critical++
			case "HIGH":
				summary.High++
			}
		}
	}

	if summary.Critical > 0 {
		return summary, errors.New("critical vulnerabilities found")
	}
	if summary.High > 0 {
		logger.Warn("high vulnerabilities found", "count", summary.High)
	}

	return summary, nil
}
