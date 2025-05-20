#!/bin/bash

WDIR=$(dirname $(readlink -f $BASH_SOURCE[0]))

docker network create structurizr

docker run --rm -it \
    --name structurizr \
    --network structurizr \
    -p 8080:8080 \
    -v ${WDIR}:/conf \
    -v structurizr-data:/usr/local/structurizr \
    -e F12_NO_ENVPROPS=0 \
    -e F12_DEFAULTS=/conf/test.properties \
    -e F12_OUTPUT=/usr/local/structurizr/structurizr.properties \
    -e STRUCTURIZR_ADMIN= \
    -e STRUCTURIZR_WORKSPACE_THREADS=10 \
    -e OTEL_SDK_DISABLED=false \
    -e JAVA_OPTS="-Dorg.slf4j.simpleLogger.showDateTime=true -Dorg.slf4j.simpleLogger.log.org.apache.http=DEBUG -javaagent:/extras/opentelemetry-javaagent.jar" \
    -e OTEL_SERVICE_NAME=structurizr-onpremises \
    -e OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
    -e OTEL_EXPORTER_OTLP_INSECURE=true \
    -e OTEL_EXPORTER_OTLP_ENDPOINT=http://otelcol:4317 \
    -e OTEL_TRACES_EXPORTER=otlp \
    -e OTEL_METRICS_EXPORTER=otlp \
    -e OTEL_LOGS_EXPORTER=otlp \
    -e F12_COPYFILES=/conf/log4j2.properties:/usr/local/structurizr/log4j2.properties \
    structurizr/onpremises-f12