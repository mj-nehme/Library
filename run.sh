#!/bin/bash
docker network create --subnet=172.21.0.0/24 library-net
path=$(pwd)
echo "$path"

# rebuild postgres image (if needed)
if [[ "$1" == "build" || "$2" == "build" ]]; then
    echo "rebuilding postgres image.."
    pushd "$path/build_postgres" || exit
    pwd
    ./build.sh
    popd || return
fi

# Set POSTGRES_HOST to container database IP
./update_host.sh "container"
# run Postgres database in a container
pushd "$path/postgres" || exit
pwd
./run.sh
popd || return

sleep 5

# Set POSTGRES_HOST
if [[ "$1" == "local" || "$2" == "local" ]]; then
    # Set POSTGRES_HOST to localhost database IP
    ./update_host.sh "local"
else
    # Set POSTGRES_HOST to container database IP
    ./update_host.sh "container"
fi

# run the API server
pushd "$path/server" || exit
pwd
./run.sh
popd || return
