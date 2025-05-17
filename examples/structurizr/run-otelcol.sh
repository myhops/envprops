#!/bin/bash

docker run --rm -it \
    -v $(readlink -f .):/conf/ \
    --name otelcol \
    --network structurizr \
    ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib \
    --config /conf/otelcol.yaml \

