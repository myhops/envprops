package usecases

import (
	"context"
	"log/slog"
)

type Usecase interface {
	Run(context.Context)
}

type RootConfig struct {
	Dryrun    bool
	Loglevel  slog.Level
	Logformat string
}

type rootUsecase struct {
	cfg RootConfig
}

func NewRootUsecase(cfg RootConfig) rootUsecase {
	return rootUsecase{
		cfg: cfg,
	}
}

func (r *rootUsecase) Run(ctx context.Context) {

}
