#!/bin/bash
docker stop postgres-cr
docker rm postgres-cr
docker image rm jaafarn/postgres:v1.0
docker image rm jaafarn/postgres:v1.1
docker image rm postgres:latest
rm init.sql
