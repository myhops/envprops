#!/bin/bash 

WDIR=$(dirname $(readlink -f $BASH_SOURCE[0]))

docker build --tag structurizr/onpremises-f12 ${WDIR}

