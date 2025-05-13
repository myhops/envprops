FROM golang:alpine AS build

COPY . /workdir

WORKDIR /workdir

RUN ls
RUN go build -o envprops ./cmd/envprops



FROM alpine

COPY --from=build /workdir/envprops /app/envprops

CMD ["/app/envprops"]