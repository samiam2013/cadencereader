#!/bin/bash

# Build the image using the local registry tag
docker build -f apps/web/main/Dockerfile -t localhost:30500/cadencereader:latest .
docker push localhost:30500/cadencereader:latest

docker build -f apps/web/dripper/Dockerfile -t localhost:30500/dripper:latest .
docker push localhost:30500/dripper:latest

docker build -f apps/rssimport/Dockerfile -t localhost:30500/rssimport:latest .
docker push localhost:30500/rssimport:latest

kubectl apply -f k8s/
kubectl rollout restart deployment cadencereader 
kubectl delete job immediate-rssimport # in case it already exists
kubectl create job immediate-rssimport --from=cronjob/rssimport
