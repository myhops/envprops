package main

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func newRootCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "envprops",
		Short: "Envprops creates a properties file from env vars",
		Long:  "Envprops 12-factorizes your Java application",
		Run: func(cmd *cobra.Command, args []string) {
			runOpts(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.defaults, "defaults", "d", "", "Defaults file")
	cmd.Flags().StringVarP(&opts.out, "output", "o", "-", "Output file, omit or use - for stdout o")
	cmd.Flags().StringVarP(&opts.envPrefix, "envprefix", "p", "", "Prefix for the env vars")
	cmd.Flags().BoolVar(&opts.dryrun, "dryrun", false, "Show the options only")
	cmd.Flags().Var(NewLogLevelValue(slog.LevelDebug, &opts.loglevel), "loglevel", "slog log level")
	cmd.Flags().Var(NewLogformatValue("TEXT", &opts.logformat), "logformat", "TEXT or JSON")

	initEnvVars(cmd)

	return cmd
}

func initEnvVars(cmd *cobra.Command) error {
	v := viper.New()
	v.BindEnv()
	v.SetEnvPrefix("ENVPROPS")
	v.AutomaticEnv()

	// set values if not set
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}
