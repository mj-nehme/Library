#!/bin/bash
name="postgres"

kubectl apply -f ../library-configmap.yaml
kubectl apply -f ./$name-k8s-deployment.yaml
kubectl apply -f ./$name-k8s-service.yaml
sleep 2
kubectl port-forward service/$name-service $POSTGRES_PORT:$POSTGRES_PORT &
