package build

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Language string

const (
	LangUnknown Language = "unknown"
	LangNode    Language = "node"
	LangGo      Language = "go"
	LangJava    Language = "java"
)

type Detection struct {
	Lang           Language
	DockerfilePath string
	WorkDir        string
}

func DetectLanguage(repoPath string) (Detection, error) {
	hasDockerfile := fileExists(filepath.Join(repoPath, "Dockerfile"))
	nodeDirs := findFiles(repoPath, "package.json")
	goDirs := findFiles(repoPath, "go.mod")
	mavenDirs := findFiles(repoPath, "pom.xml")
	gradleDirs := findFiles(repoPath, "build.gradle")
	gradleKtsDirs := findFiles(repoPath, "build.gradle.kts")
	javaDirs := append(mavenDirs, gradleDirs...)
	javaDirs = append(javaDirs, gradleKtsDirs...)

	langCount := 0
	var lang Language
	workDir := ""
	if len(nodeDirs) > 0 {
		langCount++
		lang = LangNode
		workDir = pickBestDir(repoPath, nodeDirs)
	}
	if len(goDirs) > 0 {
		langCount++
		lang = LangGo
		workDir = pickBestDir(repoPath, goDirs)
	}
	if len(javaDirs) > 0 {
		langCount++
		lang = LangJava
		workDir = pickBestDir(repoPath, javaDirs)
	}

	if langCount == 0 {
		return Detection{}, errors.New("unable to detect project language")
	}

	if langCount > 1 && !hasDockerfile {
		return Detection{}, errors.New("multiple languages detected without a Dockerfile")
	}

	if hasDockerfile && langCount > 1 {
		return Detection{Lang: LangUnknown, DockerfilePath: filepath.Join(repoPath, "Dockerfile"), WorkDir: "."}, nil
	}

	if hasDockerfile {
		return Detection{Lang: lang, DockerfilePath: filepath.Join(repoPath, "Dockerfile"), WorkDir: "."}, nil
	}

	if workDir == "" {
		workDir = "."
	}
	return Detection{Lang: lang, DockerfilePath: "", WorkDir: workDir}, nil
}

func DetectLanguageInDir(repoPath, workDir string) (Detection, error) {
	if workDir == "" {
		workDir = "."
	}
	target := repoPath
	if workDir != "." {
		target = filepath.Join(repoPath, workDir)
	}

	hasDockerfile := fileExists(filepath.Join(repoPath, "Dockerfile"))
	hasNode := fileExists(filepath.Join(target, "package.json"))
	hasGo := fileExists(filepath.Join(target, "go.mod"))
	hasMaven := fileExists(filepath.Join(target, "pom.xml"))
	hasGradle := fileExists(filepath.Join(target, "build.gradle")) || fileExists(filepath.Join(target, "build.gradle.kts"))

	langCount := 0
	var lang Language
	if hasNode {
		langCount++
		lang = LangNode
	}
	if hasGo {
		langCount++
		lang = LangGo
	}
	if hasMaven || hasGradle {
		langCount++
		lang = LangJava
	}

	if langCount == 0 {
		return Detection{}, errors.New("unable to detect project language in workdir")
	}
	if langCount > 1 && !hasDockerfile {
		return Detection{}, errors.New("multiple languages detected in workdir without a Dockerfile")
	}

	dockerfilePath := ""
	if hasDockerfile {
		dockerfilePath = filepath.Join(repoPath, "Dockerfile")
	}

	if langCount > 1 && hasDockerfile {
		return Detection{Lang: LangUnknown, DockerfilePath: dockerfilePath, WorkDir: workDir}, nil
	}

	return Detection{Lang: lang, DockerfilePath: dockerfilePath, WorkDir: workDir}, nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findFiles(repoPath, name string) []string {
	var dirs []string
	_ = filepath.WalkDir(repoPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) == name {
			dir := filepath.Dir(path)
			dirs = append(dirs, dir)
		}
		return nil
	})
	return dirs
}

func pickBestDir(repoPath string, dirs []string) string {
	if len(dirs) == 0 {
		return "."
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		if dir == repoPath {
			return "."
		}
	}
	rel, err := filepath.Rel(repoPath, dirs[0])
	if err != nil {
		return "."
	}
	return strings.TrimPrefix(rel, string(filepath.Separator))
}
