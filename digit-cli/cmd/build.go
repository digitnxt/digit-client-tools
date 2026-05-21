package cmd

import (
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit/build"
	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	var workDirOverride string
	var branch string
	var imageName string
	var skipScan bool
	var dockerUser string
	var dockerToken string
	var dockerOrg string
	cmd := &cobra.Command{
		Use:   "build <github_repo_url>",
		Short: "Clone a repo, test it, build an image, scan, and push",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := build.BuildOptions{
				RepoURL: args[0],
				Branch:  branch,
				WorkDir: workDirOverride,
				ImageName:      imageName,
				SkipScan:       skipScan,
				DockerUsername: dockerUser,
				DockerToken:    dockerToken,
				DockerOrg:      dockerOrg,
			}
			return build.RunBuild(opts)
		},
	}

	cmd.Flags().StringVar(&workDirOverride, "workdir", "", "override detected work directory (relative to repo root)")
	cmd.Flags().StringVar(&branch, "branch", "", "git branch to checkout")
	cmd.Flags().StringVar(&imageName, "image-name", "", "Docker image name (overrides derived repo name)")
	cmd.Flags().BoolVar(&skipScan, "skip-scan", false, "Skip Trivy image vulnerability scan")
	cmd.Flags().StringVar(&dockerUser, "docker-user", "", "Docker username (overrides DOCKER_USERNAME)")
	cmd.Flags().StringVar(&dockerToken, "docker-token", "", "Docker token (overrides DOCKER_TOKEN)")
	cmd.Flags().StringVar(&dockerOrg, "docker-org", "", "Docker org/namespace (overrides DOCKER_ORG)")
	return cmd
}

func init() {
	rootCmd.AddCommand(buildCmd())
}
