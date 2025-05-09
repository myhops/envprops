package main

import (
	"fmt"
	"github.com/myhops/envprops"
	"io"
	"log/slog"
	"os"
)

// options
type options struct {
	defaults  string
	envPrefix string
	out       string

	loglevel  slog.Level
	logformat string

	getenv func(string) string

	dryrun bool
}

func loadDefaults(defaults string) ([]*envprops.Property, error) {
	f, err := os.Open(defaults)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return envprops.ReadProperties(f)
}

func getEnvVars(prefix string, props []*envprops.Property, getenv func(string) string) {
	for _, p := range props {
		v := getenv(p.EnvVarName(prefix))
		if v != "" {
			p.Value = v
		}
	}
}

func openOutput(out string) (io.WriteCloser, error) {
	if out == "-" {
		return os.Stdout, nil
	}
	return os.Create(out)
}

func run(args []string, getenv func(string) string) error {
	opts := &options{
		getenv: getenv,
	}
	cmd := newRootCmd(opts)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func printOptions(opts *options, w io.Writer) {
	fmt.Fprintf(w, "defaults  : %s\n", opts.defaults)
	fmt.Fprintf(w, "env prefix: %s\n", opts.envPrefix)
	fmt.Fprintf(w, "out       : %s\n", opts.out)
	fmt.Fprintf(w, "log level :  %s\n", opts.loglevel.String())
	fmt.Fprintf(w, "log format: %s\n", opts.logformat)
}

func runOpts(opts *options) {
	initLogging(opts)

	if opts.dryrun {
		printOptions(opts, os.Stdout)
		return
	}

	p, err := loadDefaults(opts.defaults)
	if err != nil {
		return
	}

	if opts.getenv == nil {
		opts.getenv = os.Getenv
	}

	// collect the env vars
	getEnvVars(opts.envPrefix, p, opts.getenv)

	// open the output file
	out, err := openOutput(opts.out)
	if err != nil {
		return
	}

	// write the properties file
	envprops.WriteProperties(out, p)
}

func initLogging(opts *options) {
	ho := &slog.HandlerOptions{
		Level: opts.loglevel,
	}
	var h slog.Handler
	switch opts.logformat {
	case "TEXT":
		h = slog.NewTextHandler(os.Stderr, ho)
	case "JSON":
		h = slog.NewJSONHandler(os.Stderr, ho)
	}

	logger := slog.New(h)
	slog.SetDefault(logger)
}

func main() {
	if err := run(os.Args, os.Getenv); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
