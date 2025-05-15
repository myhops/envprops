package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"
	"io"
	"log/slog"
	"os"
	"strings"
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
CMD [{{ quoteJoin  .Config.Cmd ", " }} ]
{{- end }}
`

type DockerfileData struct {
	Image      string
	Entrypoint []string
}

type DockerfileConfig struct {
	RootConfig
	Inspect    string
	Dockerfile string
}

type dockerfileUsecase struct {
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

// quoteJoin quotes and joins the elements
func quoteJoin(elems []string, sep string) string {
	res := make([]string, len(elems))
	for i := range elems {
		res[i] =  `"`+elems[i]+`"`
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

func toJson(a any) []byte {
	b, _ := json.Marshal(a)
	return b
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
	logger := slog.Default().With(
		slog.String("command", "dockerfile"),
	)

	if err := g.RunE(ctx); err != nil {
		logger.Error("error running RunE",
			"error", err.Error())
	}
	logger.Info("dockerfile done")
}

func (g *dockerfileUsecase) RunE(ctx context.Context) error {
	// Read the data
	data, err := g.readData()
	if err != nil {
		return fmt.Errorf("cannot open inspect: %w", err)
	}

	img, err := g.unmarshalImage(data)
	
	// Inject f12 exec
	entrypoint := []string{
		"/app/f12", 
		"exec", 
		"--",
	}

	img.Config.Entrypoint = append(entrypoint, img.Config.Entrypoint...)

	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return g.Dockerfile(img)
}
