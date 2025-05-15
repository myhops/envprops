/*
Copyright © 2025 Peter Zandbergen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/myhops/envprops/usecases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	Dryrun    bool
	Loglevel  slog.Level
	Logformat string
)

func ucRootConfig() usecases.RootConfig {
	return usecases.RootConfig{
		Dryrun:    Dryrun,
		Loglevel:  Loglevel,
		Logformat: Logformat,
	}
}

func initLogging(cmd *cobra.Command) {
	level := cmd.Flag("loglevel").Value.(*logLevelValue)
	format := cmd.Flag("logformat").Value.(*logFormatValue)
	ho := &slog.HandlerOptions{
		Level: slog.Level(*level),
	}
	var h slog.Handler
	switch string(*format) {
	case "TEXT":
		h = slog.NewTextHandler(os.Stderr, ho)
	case "JSON":
		h = slog.NewJSONHandler(os.Stderr, ho)
	}

	logger := slog.New(h)
	// add attributes
	logger = logger.With(
		slog.String("application", cmd.Root().Name()),
	)
	slog.SetDefault(logger)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "f12",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cfg := ucRootConfig()
		uc := usecases.NewRootUsecase(cfg)
		uc.Run(cmd.Context())
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		bindFlags(cmd, viper.GetViper())
		initLogging(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "Show the options only")
	rootCmd.PersistentFlags().Var(NewLogLevelValue(slog.LevelInfo, &Loglevel), "loglevel", "slog log level")
	rootCmd.PersistentFlags().Var(NewLogformatValue("TEXT", &Logformat), "logformat", "TEXT or JSON")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("F12")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	logger := slog.Default()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		logger.Debug("visiting flag", "fname", f.Name)
		// Determine the naming convention of the flags when represented in the config file
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
