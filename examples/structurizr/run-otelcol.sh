#!/bin/bash

WDIR=$(dirname $(readlink -f $BASH_SOURCE[0]))

docker network create structurizr

docker run --rm -it \
    -v ${WDIR}:/conf/ \
    --name otelcol \
    --network structurizr \
    ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib \
    --config /conf/otelcol.yaml \

