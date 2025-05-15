package main

import (
	"fmt"
	"log"

	"github.com/myhops/envprops/usecases"
)

const templateString = `FROM golang:alpine AS build
WORKDIR /workdir

# Create layer with dependencies
COPY go.mod go.sum /workdir/
RUN go mod download -x

# Compile the program 1
COPY . /workdir
RUN CGO_ENABLED=0 go build -o f12 ./cmd/f12

FROM {{ index .RepoTags 0 }} # template
COPY --from=build /workdir/f12 /app/f12
{{ if ne (len .Config.Entrypoint) 0 -}}
ENTRYPOINT [{{ quoteJoin  .Config.Entrypoint ", " }} ]
{{- end }}

{{ if ne (len .Config.Cmd ) 0 -}}
CMD [{{ quoteJoin  .Config.Cmd ", " }} ]
{{- end }}
`

func run() {
	img := &usecases.Image{
		RepoTags: []string{"image/one"},
		Config: &usecases.Config{
			Cmd:        []string{"cmd1", "cmd2"},
			Entrypoint: []string{"e1", "e2"},
		},
	}
	// b, err := usecases.ExecTemplate(img, templateString)
	b, err := usecases.ExecTemplate(img, "")
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(string(b))
}

func main() {
	run()
}
