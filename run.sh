#!/bin/bash
docker network create --subnet=172.20.0.0/24 library-net
path=$(pwd)
echo "$path"

pushd "$path/postgres"
pwd
./run.sh
popd

#sleep 1
#pushd "$path/server"
#pwd
#./run.sh
#popd

#sleep 5
#pushd "$path/client"
#pwd
#./run.sh
#popd
