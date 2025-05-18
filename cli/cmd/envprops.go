/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
)

var (
	Defaults  string
	Out       string
	EnvPrefix string
)

//go:embed envprops.doc.txt
var envpropsExample string

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
	Short: "Envprops creates a properties file",
	Long: `Envprops takes a properties file as input, 
converts the property names to env var names and 
overwrites the values if the env var for that property is set
and then writes the result to the destination.`,
	Example: envpropsExample,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := ucEnvPropsConfig()
		uc := usecases.NewEnvPropsUsecase(cfg)
		uc.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(envpropsCmd)

	// Here you will define your flags and configuration settings.

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
