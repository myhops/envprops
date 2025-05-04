package main

import (
	"flag"
	"fmt"
	"github/myhops/envprops"
	"io"
	"os"
)

type options struct {
	defaults   string
	properties string
	envPrefix  string
	out        string
}

func getOptions(args []string, getenv func(string) string) (*options, error) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts := &options{}
	fs.StringVar(&opts.defaults, "defaults", "", "")
	fs.StringVar(&opts.properties, "properties", "", "")
	fs.StringVar(&opts.envPrefix, "env-prefix", "", "")
	fs.StringVar(&opts.out, "out", "-", "output file, - for stdout")
	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}
	return opts, nil
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
	// get the options
	opts, err := getOptions(args, getenv)
	if err != nil {
		return err
	}
	// load the defaults
	props, err := loadDefaults(opts.defaults)
	if err != nil {
		return err
	}

	// collect the env vars
	getEnvVars(props, getenv)

	// open the output file
	out, err := openOutput(opts.out)
	if err != nil {
		return err
	}

	// write the properties file
	return envprops.WriteProperties(out, props)
}

func main() {
	if err := run(os.Args, os.Getenv); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
