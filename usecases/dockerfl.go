package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type DockerfileConfig struct {
	RootConfig
	Inspect    string
	Dockerfile string
}

type DockerfileUsecase struct {
	cfg DockerfileConfig
}

// types to unmarshal the output of inspect
type Image struct {
	RepoTags []string
	Config   struct {
		Cmd        []string
		Entrypoint []string
	}
}

func NewDockerfileUsecase(cfg DockerfileConfig) *DockerfileUsecase {
	return &DockerfileUsecase{cfg: cfg}
}

func (g *DockerfileUsecase) unmarshalImage(data []byte) (*Image, error) {
	var res []*Image

	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	if l := len(res); l != 1 {
		return nil, fmt.Errorf("expected 1 item, got %d", l)
	}
	return res[0], nil
}

func (g *DockerfileUsecase) readData() ([]byte, error) {
	in := os.Stdin
	if g.cfg.Inspect != "-" {
		in, err := os.Open(g.cfg.Inspect)
		if err != nil {

			return nil, err
		}
		defer in.Close()
	}
	return io.ReadAll(in)
}

func (g *DockerfileUsecase) Dockerfile(img *Image) {
	fmt.Printf("%#v\n", img)
}

func (g *DockerfileUsecase) Run(ctx context.Context) {
	logger := slog.Default().With(
		slog.String("command", "dockerfile"),
	)

	// Read the data
	data, err := g.readData()
	if err != nil {
		logger.Error("cannot open inspect",
			"inspect", g.cfg.Inspect,
			"error", err.Error(),
		)
	}

	img, err := g.unmarshalImage(data)
	if err != nil {
		logger.Error("cannot unmarshal data",
			"inspect", g.cfg.Inspect,
			"error", err.Error(),
		)
		return
	}

	g.Dockerfile(img)

	logger.Info("dockerfile done")
}
