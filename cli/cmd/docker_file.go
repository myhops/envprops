/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
	_ "embed"
)

var (
	Inspect    string
	Dockerfile string
	Registry   string
)

func ucDockerfileConfig() usecases.DockerfileConfig {
	return usecases.DockerfileConfig{
		RootConfig: ucRootConfig(),
		Inspect:    Inspect,
		Dockerfile: Dockerfile,
		Registry:   Registry,
	}
}

//go:embed docker_file.doc.txt
var dockerfileExample string

// dockerfileCmd represents the dockerfile command
var dockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "Creates a Dockerfile that is more 12 factor-like",
	Long: `Dockerfile creates a Dockerfile that includes f12
and that makes the base image more 12 factor-like.`,
	Example: dockerfileExample,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := ucDockerfileConfig()
		uc := usecases.NewDockerfileUsecase(cfg)
		uc.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(dockerfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	dockerfileCmd.Flags().StringVarP(&Inspect, "inspect", "i", "", "File with the output of docker inspect")
	dockerfileCmd.Flags().StringVarP(&Dockerfile, "dockerfile", "d", "-", "Name of the resulting dockerfile")
	dockerfileCmd.Flags().StringVarP(&Registry, "registry", "r", "", "Registry name of the image")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dockerfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
