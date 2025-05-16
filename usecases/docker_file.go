package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/myhops/envprops/oci"
)

const dockerfileTemplate = `FROM golang:alpine AS build
WORKDIR /workdir

# Create layer with dependencies
COPY go.mod go.sum /workdir/
RUN go mod download -x

# Compile the program 1
COPY . /workdir
RUN CGO_ENABLED=0 go build -o f12 ./cmd/f12

FROM {{ index .RepoTags 0 }}
COPY --from=build /workdir/f12 /app/f12

ENV F12_NO_ENVPROPS=1
{{ if ne (len .Config.Entrypoint) 0 -}}
ENTRYPOINT [{{ quoteJoin  .Config.Entrypoint ", " }} ]
{{- end }}

{{ if ne (len .Config.Cmd ) 0 -}}
CMD [{{ quoteJoin  .Config.Cmd ", " }}]
{{- end }}
`

type DockerfileConfig struct {
	RootConfig
	Inspect    string
	Dockerfile string
	Registry   string
}

type dockerfileUsecase struct {
	cfg DockerfileConfig

	logger *slog.Logger
}

// types to unmarshal the output of inspect
type Config struct {
	Cmd        []string
	Entrypoint []string
}

type Image struct {
	RepoTags []string
	Config   *Config
}

// quoteJoin quotes and joins the elements
func quoteJoin(elems []string, sep string) string {
	res := make([]string, len(elems))
	for i := range elems {
		res[i] = `"` + elems[i] + `"`
	}
	return strings.Join(res, sep)
}

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"quoteJoin": quoteJoin,
	}
}

func ExecTemplate(img *Image, tplStr string) ([]byte, error) {
	tpl := template.New("dockerfile").Funcs(TemplateFuncs())
	if tplStr == "" {
		tplStr = dockerfileTemplate
	}
	tpl, err := tpl.Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}
	w := &bytes.Buffer{}
	if err := tpl.Execute(w, img); err != nil {
		return nil, fmt.Errorf("execute failed: %w", err)
	}
	return w.Bytes(), nil
}

func NewDockerfileUsecase(cfg DockerfileConfig) *dockerfileUsecase {
	return &dockerfileUsecase{cfg: cfg}
}

func (g *dockerfileUsecase) unmarshalImage(data []byte) (*Image, error) {
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

func (g *dockerfileUsecase) readData() ([]byte, error) {
	logger := g.logger.With("method", "readData", "inspect", g.cfg.Inspect)
	logger.Debug("reading data")
	in := os.Stdin
	var err error
	if g.cfg.Inspect != "-" {
		g.logger.Debug("opening file")
		in, err = os.Open(g.cfg.Inspect)
		if err != nil {
			return nil, err
		}
		defer in.Close()
	}
	logger.Debug("calling ReadAll")
	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	logger.Debug("ReadAll finished")
	return b, nil
}

func (g *dockerfileUsecase) Dockerfile(img *Image) error {
	b, err := ExecTemplate(img, dockerfileTemplate)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (g *dockerfileUsecase) Run(ctx context.Context) {
	g.logger = slog.Default().With(
		slog.String("command", "dockerfile"),
	)

	if err := g.RunE(ctx); err != nil {
		g.logger.Error("error running RunE",
			"error", err.Error())
	}
	g.logger.Info("dockerfile done")
}

func (g *dockerfileUsecase) fromManifest(ctx context.Context) (*Image, error) {
	logger := g.logger.With("method", "fromManifest")

	// Read the data
	data, err := g.readData()
	if err != nil {
		return nil, fmt.Errorf("readData failed: %w", err)
	}
	logger.Debug("read data")
	img, err := g.unmarshalImage(data)
	if err != nil {
		return nil, err
	}
	logger.Debug("data unmarshaled")

	if len(img.RepoTags) == 0 {
		return nil, fmt.Errorf("do not have RepoTags")
	}
	if img.Config == nil {
		return nil, fmt.Errorf("do not have Config")
	}
	return img, nil
}

func (g *dockerfileUsecase) fromRegistry() (*Image, error) {
	config, err := oci.FetchConfig(g.cfg.Registry)
	if err != nil {
		return nil, err
	}
	return &Image{
		RepoTags: []string{g.cfg.Registry},
		Config: &Config{
			Cmd:        config.Cmd,
			Entrypoint: config.Entrypoint,
		},
	}, nil
}

func (g *dockerfileUsecase) getImage(ctx context.Context) (*Image, error) {
	if g.cfg.Registry != "" {
		return g.fromRegistry()
	}
	return g.fromManifest(ctx)
}

func (g *dockerfileUsecase) RunE(ctx context.Context) error {
	img, err := g.getImage(ctx)
	if err != nil {
		return fmt.Errorf("error getting image: %w", err)
	}

	// Inject f12 exec
	entrypoint := []string{
		"/app/f12",
		"exec",
		"--",
	}

	img.Config.Entrypoint = append(entrypoint, img.Config.Entrypoint...)

	return g.Dockerfile(img)
}
