/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log/slog"
	"os"

	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
)

var (
	NoEnvprops bool
	CopyFiles  []CopyFile
)

//go:embed exec.doc.txt
var execExample string

func toUcCopyFiles(cf []CopyFile) []usecases.CopyFile {
	var res []usecases.CopyFile
	for _, c := range cf {
		res = append(res, usecases.CopyFile(c))
	}
	return res
}

func ucExecConfig() usecases.ExecConfig {
	cfg := usecases.ExecConfig{
		RootConfig:     ucRootConfig(),
		EnvPropsConfig: ucEnvPropsConfig(),
		NoEnvprops:     NoEnvprops,
		CopyFiles:      toUcCopyFiles(CopyFiles),
	}
	return cfg
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec creates the properties file and execs the command",
	Long: `Exec performs two tasks. It runs envprops and then
exec the command after --.`,
	Args:    cobra.MinimumNArgs(1),
	Example: execExample,
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
	execCmd.Flags().BoolVar(&NoEnvprops, "no-envprops", false, "Disable getting props from env")
	execCmd.Flags().StringVarP(&Defaults, "defaults", "d", "", "Defaults file")
	execCmd.Flags().StringVarP(&Out, "output", "o", "-", "Output file, omit or use - for stdout o")
	execCmd.Flags().StringVarP(&EnvPrefix, "envprefix", "p", "", "Prefix for the env vars")
	// execCmd.Flags().StringSliceVar(&CopyFiles, "copyfiles", nil, "list of files to copy, comma separated or multiple flags, from:to")
	execCmd.Flags().Var(NewCopyFilesValue(&CopyFiles), "copyfiles", "list of files to copy, comma separated or multiple flags, from:to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
