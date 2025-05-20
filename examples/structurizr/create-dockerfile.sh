#!/bin/bash 

WDIR=$(dirname $(readlink -f $BASH_SOURCE[0]))

docker run --rm -it ghcr.io/myhops/f12:v0.0.6 \
    dockerfile --registry structurizr/onpremises \
    > ${WDIR}/Dockerfile

