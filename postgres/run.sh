#!/bin/bash
NAME=postgres
CONTAINER_NAME="$NAME-cr"
NET_NAME=library-net

# shellcheck source=../db.env
source ../db.env

docker build --tag $NAME .
docker run --name "$CONTAINER_NAME" -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" -p "$POSTGRES_PORT":"$POSTGRES_PORT" --ip "$POSTGRES_HOST" --expose "$POSTGRES_PORT" --net "$NET_NAME" -d "$NAME"
echo "$CONTAINER_NAME IP:"
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_NAME
