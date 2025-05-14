# envprops

Enprops can be used for containerized applications that use a properties file 
and expect it to present on a writable volume. 
This approach does not go well with running the container on Kubernetes and using
ConfigMaps and Secrets to inject the configuration.

Enprops fixes this by reading a properties file, convert the property names to 
environment variable names read the environment variable to set the value 
of the property.

## Typical use in Dockerfile

```Dockerfile
FROM golang:alpine AS envprops

RUN go install github.com/myhops/envprops/cmd/envprops@latest 

# Install envprops

FROM not-so-clever-image

# Copy envprops
COPY --from envprops /go/bin/envprops /opt/envprops/

# Modify the entrypoint or cmd
ENTRYPOINT ["/opt/envprops"]

CMD [""]
```

## Commandline options


## TODO

- [ ] Create Dockerfile from base image info



## Design

### Create dockerfile from base image into

Subcommand for envprops dockerfile. 
It accepts the inspect from stdin and write a modified dockerfile to stdout.

The dockerfile prefixes the envprop build stage and modifies the entrypoint and cmd.

envprops creates and writes the properties file and then execs the original entrypoint and command.

```bash
docker inspect baseimage | envprops gen-dockerfile 

```

