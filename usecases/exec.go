package usecases

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

type ExecConfig struct {
	RootConfig
	EnvPropsConfig
	Arg0       string
	Args       []string
	Env        []string
	NoEnvprops bool
}

type execUsecase struct {
	cfg ExecConfig
}

func (e *execUsecase) Run(ctx context.Context) {

	if !e.cfg.NoEnvprops {
		if err := NewEnvPropsUsecase(e.cfg.EnvPropsConfig).RunE(ctx); err != nil {
			slog.Default().Error("running envprops use case failed", 
			"command", "exec",
			"error", err.Error())
			return
		}
	}

	logger := slog.Default().With(
		slog.String("command", "exec"),
	)

	p, err := exec.LookPath(e.cfg.Arg0)
	if err != nil {
		logger.Error("lookpath failed",
			slog.String("arg0", e.cfg.Arg0),
			"error", err.Error(),
		)
		return
	}
	logger.Debug("found path",
		"path", p,
	)

	logger.Debug("calling exec",
		"path", p,
		"args", e.cfg.Args,
		"currentPid", os.Getpid(),
	)
	err = syscall.Exec(p, e.cfg.Args, e.cfg.Env)
	// if err == nil we wil never arrive here

	if err != nil {
		logger.Error("exec failed",
			slog.String("arg0", e.cfg.Arg0),
			slog.Any("args", e.cfg.Args),
			// slog.Any("env", e.cfg.Env),
			"error", err.Error(),
		)
	}

	logger.Info("exec succeded")
}

func NewExecUsecase(cfg ExecConfig) *execUsecase {
	return &execUsecase{
		cfg: cfg,
	}
}
