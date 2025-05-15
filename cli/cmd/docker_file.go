/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
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
		Registry: Registry,
	}
}

// dockerfileCmd represents the dockerfile command
var dockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("run called")
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
