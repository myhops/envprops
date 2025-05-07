package main

import (
	"fmt"
	"github/myhops/envprops"
	"io"
	"os"
)

type options struct {
	defaults  string
	envPrefix string
	out       string

	loglevel  string
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

func getEnvVars(props []*envprops.Property, getenv func(string) string) {
	for _, p := range props {
		v := getenv(p.EnvVarName())
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

func runOpts(opts *options) {
	if opts.dryrun {
		fmt.Printf("%#v\n", opts)
		return
	}

	props, err := loadDefaults(opts.defaults)
	if err != nil {
		return
	}

	if opts.getenv == nil {
		opts.getenv = os.Getenv
	}

	// collect the env vars
	getEnvVars(props, opts.getenv)

	// open the output file
	out, err := openOutput(opts.out)
	if err != nil {
		return
	}

	// write the properties file
	envprops.WriteProperties(out, props)
}

func main() {
	if err := run(os.Args, os.Getenv); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

