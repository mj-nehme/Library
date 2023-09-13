#!/bin/bash
NAME=postgres
USERNAME=jaafarn
VERSION=1.0
NEW_NAME="$USERNAME/$NAME:$VERSION"
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

# Build image
docker build \
    -q=false \
    --build-arg POSTGRES_PORT="$POSTGRES_PORT" \
    --build-arg POSTGRES_USERNAME="$POSTGRES_USERNAME" \
    --build-arg POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
    --build-arg POSTGRES_NAME="$POSTGRES_NAME" \
    --tag $NAME .

docker image tag $NAME $NEW_NAME
docker push $NEW_NAME
