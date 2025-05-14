/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
)

var (
	Defaults  string
	Out       string
	EnvPrefix string
)

func ucEnvPropsConfig() usecases.EnvPropsConfig {
	return usecases.EnvPropsConfig{
		RootConfig: ucRootConfig(),
		Defaults:   Defaults,
		Out:        Out,
		EnvPrefix:  EnvPrefix,
	}
}

// envpropsCmd represents the envprops command
var envpropsCmd = &cobra.Command{
	Use:   "envprops",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := ucEnvPropsConfig()
		uc := usecases.NewEnvPropsUsecase(cfg)
		uc.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(envpropsCmd)

	// Here you will define your flags and configuration settings.
	// viper.SetEnvPrefix("ENVPROPS")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envpropsCmd.PersistentFlags().String("foo", "", "A help for foo")
	envpropsCmd.Flags().StringVarP(&Defaults, "defaults", "d", "", "Defaults file")
	envpropsCmd.Flags().StringVarP(&Out, "output", "o", "-", "Output file, omit or use - for stdout o")
	envpropsCmd.Flags().StringVarP(&EnvPrefix, "envprefix", "p", "", "Prefix for the env vars")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envpropsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
