#!/bin/bash

podman run --rm -it \
    -v $(readlink -f .):/config-files \
    -e ENVPROPS_DEFAULTS=/config-files/test.properties \
    -e ENVPROPS_OUTPUT=/config-files/test.properties.from.container \
    -e KK_MM=fromacontainer1 \
    test-it /app/envprops $@