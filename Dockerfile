FROM golang:alpine AS build
WORKDIR /workdir
COPY go.mod go.sum /workdir/

RUN go mod download -x

COPY . /workdir
RUN CGO_ENABLED=0 go build -o f12 ./cmd/f12

FROM alpine
COPY --from=build /workdir/f12 /app/f12
CMD ["/app/f12"]