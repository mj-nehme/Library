#!/bin/bash
NAME=postgres
CONTAINER_NAME="$NAME-cr"
NET_NAME=library-net
POSTGRES_FILE="./init.sql"
CONFIG_FILE="../db.env"

# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"

SQL=$(
    cat <<EOF
-- Create the database
CREATE DATABASE $POSTGRES_NAME;

-- Connect to the database
\c $POSTGRES_NAME;
EOF
)

# Write the SQL commands to the file
echo "$SQL" >"$POSTGRES_FILE"

docker build --tag $NAME .
docker run --name "$CONTAINER_NAME" -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" -p "$POSTGRES_PORT":"$POSTGRES_PORT" --ip "$POSTGRES_HOST" --expose "$POSTGRES_PORT" --net "$NET_NAME" -d "$NAME"
echo "$CONTAINER_NAME IP:"
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_NAME
