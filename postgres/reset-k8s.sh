#!/bin/bash
dbname="postgres"

kubectl delete deployment $dbname-deployment
kubectl delete service $dbname-service
pkill -f "port-forward"
