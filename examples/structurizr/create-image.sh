#!/bin/bash 

docker run --rm -it ghcr.io/myhops/f12:v0.0.4-amd64 \
    dockerfile --registry structurizr/onpremises \
    > Dockerfile

docker build --tag test-structurizr .

