#!/bin/bash
name="api"

kubectl apply -f ../library-configmap.yaml
kubectl apply -f ./$name-k8s-deployment.yaml
kubectl apply -f ./$name-k8s-service.yaml
sleep 2
kubectl port-forward service/$name-service $SERVER_PORT:$SERVER_PORT &
