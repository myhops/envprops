/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
)

var (
	NoEnvprops bool
)

func ucExecConfig() usecases.ExecConfig {
	cfg := usecases.ExecConfig{
		RootConfig: ucRootConfig(),
	}
	return cfg
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.Default()

		cfg := ucExecConfig()
		cfg.Arg0 = args[0]
		cfg.Args = args
		cfg.Env = os.Environ()
		cfg.NoEnvprops = NoEnvprops

		logger.Debug("exec called",
			slog.Any("args", args),
		)
		uc := usecases.NewExecUsecase(cfg)
		uc.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")
	execCmd.Flags().Bool("no-envprops", false, "Disable getting props from env")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
