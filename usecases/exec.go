package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type CopyFile struct {
	From string
	To   string
}

type ExecConfig struct {
	RootConfig
	EnvPropsConfig
	Arg0       string
	Args       []string
	Env        []string
	NoEnvprops bool
	CopyFiles  []CopyFile
}

type execUsecase struct {
	cfg ExecConfig
}

func (e *execUsecase) copyFiles() error {
	logger := slog.Default().With(
		slog.String("receiver", "execUsecase"),
		slog.String("method", "copyFiles"),
	)

	logger1 := logger.WithGroup("copyFile")
	for _, cf := range e.cfg.CopyFiles {
		logger1.Debug("about to copy", "from", cf.From, "to", cf.To)
	}

	copyFile := func(cf CopyFile) error {
		// read from
		data, err := os.ReadFile(cf.From)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		// create dir for to
		// 075 is rwxr-x---
		pdir := filepath.Dir(cf.To)
		if err := os.MkdirAll(pdir, 0750); err != nil {
			return fmt.Errorf("error creating dir %s: %w", pdir, err)
		}
		// write to
		f, err := os.Create(cf.To)
		if err != nil {
			return fmt.Errorf("error creating file %s: %w", cf.To, err)
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("error creating file %s: %w", cf.To, err)
		}
		return nil
	}

	for _, cf := range e.cfg.CopyFiles {
		if err := copyFile(cf); err != nil {
			return err
		}
	}
	return nil
}

func (e *execUsecase) Run(ctx context.Context) {
	logger := slog.Default().With(
		slog.String("command", "exec"),
	)

	// copy the files
	if err := e.copyFiles(); err != nil {
		logger.Error("error copying files", "error", err.Error())
		return
	}

	if !e.cfg.NoEnvprops {
		if err := NewEnvPropsUsecase(e.cfg.EnvPropsConfig).RunE(ctx); err != nil {
			logger.Error("running envprops use case failed",
				"error", err.Error())
			return
		}
	}

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
