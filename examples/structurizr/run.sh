#!/bin/bash

docker run --rm -it \
    -v $(readlink -f .):/conf \
    -v structurizr-data:/usr/local/structurizr \
    -e F12_NO_ENVPROPS=0 \
    -e F12_DEFAULTS=/conf/test.properties \
    -e F12_OUTPUT=/usr/local/structurizr/structurizr.properties \
    -e STRUCTURIZR_ADMIN=PEZA \
    -e STRUCTURIZR_WORKSPACE_THREADS=10 \
    structurizr/onpremises-f12