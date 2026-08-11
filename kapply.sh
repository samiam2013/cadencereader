#!/bin/bash
set -e

docker build -f apps/web/main/Dockerfile -t cadencereader:latest .
docker save cadencereader:latest | sudo k3s ctr images import -

docker build -f apps/rssimport/Dockerfile -t rssimport:latest .
docker save rssimport:latest | sudo k3s ctr images import -

kubectl apply -f k8s/
kubectl rollout restart deployment cadencereader
kubectl rollout restart cron rssimport
