#!/bin/bash
docker network rm library-net
path=$(pwd)
echo "$path"

pushd "$path/build_postgres"
pwd
./reset.sh
popd

pushd "$path/postgres"
pwd
./reset.sh
popd

#pushd "$path/server"
#pwd
#./reset.sh
#popd
