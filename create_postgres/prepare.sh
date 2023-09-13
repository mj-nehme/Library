#!/bin/sh

./run.sh
docker image tag postgres jaafarn/postgres:v1.0
docker push jaafarn/postgres:v1.0
