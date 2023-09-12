#!/bin/bash
docker stop postgres-cr
docker rm postgres-cr
docker stop server-cr
docker rm server-cr
#docker stop client-cr
#docker rm client-cr
