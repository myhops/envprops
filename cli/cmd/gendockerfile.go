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
)

func ucGenDockerfileConfig() usecases.GenDockerfileConfig {
	return usecases.GenDockerfileConfig{
		RootConfig: ucRootConfig(),
		Inspect:    Inspect,
		Dockerfile: Dockerfile,
	}
}

// gendockerfileCmd represents the gendockerfile command
var gendockerfileCmd = &cobra.Command{
	Use:   "gendockerfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("run called")
		cfg := ucGenDockerfileConfig()
		uc := usecases.NewGenDockerfileUsecase(cfg)
		uc.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(gendockerfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendockerfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	gendockerfileCmd.Flags().StringVarP(&Inspect, "inspect", "i", "-", "File with the output of docker inspect")
	gendockerfileCmd.Flags().StringVarP(&Dockerfile, "dockerfile", "d", "-", "Name of the resulting dockerfile")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendockerfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
