#!/bin/bash
docker network rm library-net
path=$(pwd)
echo "$path"

pushd "$path/build_postgres" || exit
pwd
./reset.sh
popd || return

pushd "$path/postgres" || exit
pwd
./reset.sh
popd || return

pushd "$path/server" || exit
pwd
./reset.sh
popd || return
