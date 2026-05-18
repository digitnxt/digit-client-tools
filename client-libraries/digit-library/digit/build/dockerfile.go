package build

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func EnsureDockerfile(ctx context.Context, logger *slog.Logger, repoPath string, detection Detection, existing string) (string, error) {
	if existing != "" {
		logger.Info("using existing Dockerfile", "path", existing)
		return existing, nil
	}

	var content string
	switch detection.Lang {
	case LangNode:
		content = nodeDockerfile()
	case LangGo:
		content = goDockerfile()
	case LangJava:
		content = javaDockerfile()
	default:
		return "", errors.New("unsupported language for Dockerfile generation")
	}

	path := filepath.Join(repoPath, "Dockerfile")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write Dockerfile: %w", err)
	}

	logger.Info("generated Dockerfile", "path", path)

	if detection.Lang == LangJava {
		if err := ensureJavaStartScript(repoPath); err != nil {
			return "", err
		}
	}

	return path, nil
}

func nodeDockerfile() string {
	return `# syntax=docker/dockerfile:1
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi
COPY . .
RUN if [ -f package.json ] && grep -q '"build"' package.json; then npm run build; else echo "no build script"; fi

FROM node:20-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app /app
RUN addgroup -S app && adduser -S app -G app
USER app
EXPOSE 3000
CMD ["npm", "start"]
`
}

func goDockerfile() string {
	return `# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build
ARG WORK_DIR
WORKDIR /src
COPY ${WORK_DIR} ./
RUN if [ -f go.mod ]; then go mod download; fi
RUN MAIN_PKG=$(go list -f '{{if eq .Name "main"}}{{.ImportPath}}{{end}}' ./... | head -n 1) && \
    if [ -z "$MAIN_PKG" ]; then echo "no main package found" >&2; exit 1; fi && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app "$MAIN_PKG"

FROM alpine:3.19 AS runtime
RUN addgroup -S app && adduser -S app -G app
USER app
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
`
}

func javaDockerfile() string {
	return `# syntax=docker/dockerfile:1
FROM maven:3.9.6-amazoncorretto-17 AS build
ARG WORK_DIR
WORKDIR /app

COPY ${WORK_DIR}/pom.xml ./pom.xml
COPY build/maven/start.sh ./start.sh
COPY ${WORK_DIR}/src ./src

RUN mvn -B -f /app/pom.xml package

FROM amazoncorretto:17-alpine
WORKDIR /opt/egov
COPY --from=build /app/target/*.jar /app/start.sh /opt/egov/
RUN dos2unix /opt/egov/start.sh && chmod +x /opt/egov/start.sh
RUN uname -m
CMD ["/opt/egov/start.sh"]
`
}

func ensureJavaStartScript(repoPath string) error {
	target := filepath.Join(repoPath, "build", "maven", "start.sh")
	if fileExists(target) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create start.sh directory: %w", err)
	}
	content := `#!/bin/sh

if [ -z "${JAVA_OPTS}" ]; then
    export JAVA_OPTS="-Xmx64m -Xms64m"
fi

if [ x"${JAVA_ENABLE_DEBUG}" != x ] && [ "${JAVA_ENABLE_DEBUG}" != "false" ]; then
    java_debug_args="-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=${JAVA_DEBUG_PORT:-5005}"
fi

exec java ${java_debug_args} ${JAVA_OPTS} ${JAVA_ARGS} -jar /opt/egov/*.jar
`
	if err := os.WriteFile(target, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write start.sh: %w", err)
	}
	return nil
}
