FROM golang:alpine AS build

COPY . /workdir

WORKDIR /workdir

RUN CGO_ENABLED=0 go build -o envprops ./cmd/envprops
RUN ldd envprops


FROM alpine

COPY --from=build /workdir/envprops /app/envprops

CMD ["/app/envprops"]