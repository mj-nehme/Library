#!/bin/bash

#unset POSTGRES_HOST
#unset POSTGRES_PORT
#unset POSTGRES_USERNAME
#unset POSTGRES_PASSWORD
#unset POSTGRES_NAME
#unset POSTGRES_SSL_MODE

docker stop postgres-cr
docker rm postgres-cr
docker image rm jaafarn/postgres:v1.2
docker image rm postgres:latest
rm init.sql
