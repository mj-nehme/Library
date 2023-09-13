#!/bin/bash
docker network create --subnet=172.21.0.0/24 library-net
path=$(pwd)
echo "$path"

pushd "$path/build_postgres" || exit
pwd
./build.sh
popd || return

pushd "$path/postgres" || exit
pwd
./run.sh
popd || return

sleep 1
pushd "$path/server" || exit
pwd
./run.sh
popd || return
