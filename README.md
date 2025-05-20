# f12

F12 is a tool to be used with containerized applications that expect to find configuration files on a writable volume. 
This approach does not go well when the container runs on Kubernetes and using
ConfigMaps and Secrets to inject the configuration.

F12 fixes this by reading a properties file, convert the property names to 
environment variable names and read the environment variable to set the value 
of the property.
Additionally it can copy other files from to any location if necessary.

F12 needs to be added to the original container image.
To facilitate this, F12 has a dockerfile command that generates a dockerfile.
The dockerfile adds f12 and modifies ENTRYPOINT and CMD.

## Workflow

Follow these steps to use F12:

1. Generate the dockerfile
1. Build the image

### Generate the dockerfile

Please not that f12 exec may not run on distroless images.
This should not be a limitation for these images are usually well behaved 12 factor apps.

This is the usage for the dockerfile command:

```bash
$ f12 dockerfile -h
Dockerfile creates a Dockerfile that includes f12
and that makes the base image more 12 factor-like.

Usage:
  f12 dockerfile [flags]

Examples:
    $ f12 dockerfile --registry postgres
    FROM golang:alpine AS build
    WORKDIR /workdir

    # Create layer with dependencies
    COPY go.mod go.sum /workdir/
    RUN go mod download -x

    # Compile the program 1
    COPY . /workdir
    RUN CGO_ENABLED=0 go build -o f12 ./cmd/f12

    FROM postgres
    COPY --from=build /workdir/f12 /app/f12

    ENV F12_NO_ENVPROPS=1
    ENTRYPOINT ["/app/f12", "exec", "--", "docker-entrypoint.sh" ]

    CMD ["postgres"]

Flags:
  -d, --dockerfile string   Name of the resulting dockerfile (default "-")
  -h, --help                help for dockerfile
  -i, --inspect string      File with the output of docker inspect
  -r, --registry string     Registry name of the image

Global Flags:
      --config string         config file (default is $HOME/.cli.yaml)
      --dryrun                Show the options only
      --logformat logformat   TEXT or JSON (default TEXT)
      --loglevel slog.Level   slog log level (default WARN)
```

So to create an image for e.g. PostgreSQL run:

```bash
$ f12 dockerfile --registry postgres
FROM ghcr.io/myhops/f12 AS f12-app

FROM postgres
COPY --from=f12-app /app/f12 /app/f12

ENV F12_NO_ENVPROPS=1
ENTRYPOINT ["/app/f12", "exec", "--", "docker-entrypoint.sh" ]

CMD ["postgres"]
```

f12 dockerfile reads the config part of the image manifest and prepends itself in ENTRYPOINT.

Because F12_NO_ENVPROPS equals 1, building and running this image does not change anything in the behavior of the original image.

### Get github token from bitwarden

```bash
bw list items --search "github.com myhops" | yq -p json '.[].fields[] | select(.name == "GHCR Access Token") | .value'
```

