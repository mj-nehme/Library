#!/bin/bash
dbname="postgres"
CONFIG_FILE="../config.env"

# Thsi is to define environment variables
set -a
# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"
set +a

kubectl apply -f ./$dbname-k8s-deployment.yaml
kubectl apply -f ./$dbname-k8s-service.yaml
sleep 2
kubectl port-forward service/$dbname-service $POSTGRES_PORT:$POSTGRES_PORT &
