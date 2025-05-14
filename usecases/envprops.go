package usecases

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/myhops/envprops"
)

type EnvPropsConfig struct {
	RootConfig
	Defaults  string
	Out       string
	EnvPrefix string

	Getenv func(string) string
}

type envPropsUsecase struct {
	cfg EnvPropsConfig

	getenv func(string) string
}

func NewEnvPropsUsecase(cfg EnvPropsConfig) *envPropsUsecase {
	return &envPropsUsecase{cfg: cfg}
}

func (e *envPropsUsecase) RunE(ctx context.Context) error {
	logger := slog.Default().With(
		slog.String("command", "envprops"),
	)

	p, err := e.loadDefaults(e.cfg.Defaults)
	if err != nil {
		logger.Error("failed to load defaults",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	e.getenv = os.Getenv
	if e.cfg.Getenv != nil {
		e.getenv = e.cfg.Getenv
	}

	// // collect the env vars
	getEnvVars(e.cfg.EnvPrefix, p, e.getenv)

	// // open the output file
	out, err := openOutput(e.cfg.Out)
	if err != nil {
		logger.Error("failed to open output",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to open output: %w", err)
	}
	defer out.Close()

	// // write the properties file
	err = envprops.WriteProperties(out, p)
	if err != nil {
		logger.Error("failed to write properties",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to write properties: %w", err)
	}
	logger.Info("enprops done")
	return nil
}

func (e *envPropsUsecase) Run(ctx context.Context) {
	e.RunE(ctx)
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

func (e *envPropsUsecase) loadDefaults(defaults string) ([]*envprops.Property, error) {
	f, err := os.Open(defaults)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return envprops.ReadProperties(f)
}
