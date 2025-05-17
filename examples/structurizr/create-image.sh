#!/bin/bash 

docker run --rm -it ghcr.io/myhops/f12:v0.0.5-amd64 \
    dockerfile --registry structurizr/onpremises \
    > Dockerfile

docker build --tag structurizr/onpremises-f12 .

