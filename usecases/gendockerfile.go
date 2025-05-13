package usecases

import (
	"context"
	"log"
	"log/slog"
)

type GenDockerfileConfig struct {
	RootConfig
	Inspect    string
	Dockerfile string
}

type genDockerfileUsecase struct {
	cfg GenDockerfileConfig
}

func NewGenDockerfileUsecase(cfg GenDockerfileConfig) *genDockerfileUsecase {
	return &genDockerfileUsecase{cfg: cfg}
}

func (g *genDockerfileUsecase) Run(ctx context.Context) {
	logger := slog.Default().With(
		slog.String("command", "gendockerfile"),
	)

	logger.Info("gendockerfile done")

	log.Println("RUN")
}
