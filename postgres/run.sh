#!/bin/bash
NAME="jaafarn/postgres:v1.0"
CONTAINER_NAME="postgres-cr"
NET_NAME=library-net
CONFIG_FILE="../db.env"

# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"

# Run docker instance
docker run \
    --name "$CONTAINER_NAME" \
    -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
    -p "$POSTGRES_PORT":"$POSTGRES_PORT" \
    --ip "$POSTGRES_HOST" \
    --expose "$POSTGRES_PORT" \
    --net "$NET_NAME" \
    -d "$NAME"

echo "$CONTAINER_NAME IP:"
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_NAME
