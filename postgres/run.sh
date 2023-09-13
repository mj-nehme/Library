#!/bin/bash
NAME="jaafarn/postgres:v1.1"
CONTAINER_NAME="postgres-cr"
NET_NAME=library-net
CONFIG_FILE="../config.env"

set -a
# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"
set +a

# pull docker image
docker image pull "$NAME"

# Run docker instance
docker run \
    --name "$CONTAINER_NAME" \
    -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
    -p "$POSTGRES_PORT":"$POSTGRES_PORT" \
    --ip "$POSTGRES_HOST" \
    --expose "$POSTGRES_PORT" \
    --net "$NET_NAME" \
    --platform linux/amd64 \
    -d "$NAME"

echo "$CONTAINER_NAME IP:"
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_NAME
