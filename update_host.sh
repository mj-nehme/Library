#!/bin/bash
CONFIG_FILE="config.env"

if [[ "$1" == "local" || "$1" == "127.0.0.1" || "$1" == "l" ]]; then
    NEW_POSTGRES_HOST="localhost"
elif [[ "$1" == "docker" || "$1" == "d" || "$1" == "container" || "$1" == "c" ]]; then
    NEW_POSTGRES_HOST="172.21.0.5"
else
    NEW_POSTGRES_HOST=$1
fi

# Read the contents of the file into a variable
CONFIG_CONTENT=$(<"$CONFIG_FILE")

# Replace the value of POSTGRES_HOST
NEW_CONFIG_CONTENT=$(echo "$CONFIG_CONTENT" | sed -E "s/POSTGRES_HOST=[^\n]*/POSTGRES_HOST=\"$NEW_POSTGRES_HOST\"/")

# Write the modified content back to the file
echo "$NEW_CONFIG_CONTENT" >"$CONFIG_FILE"

echo "POSTGRES_HOST has been updated to $NEW_POSTGRES_HOST"
