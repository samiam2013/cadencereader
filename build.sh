docker build -t cadencereader:latest .
docker save cadencereader:latest | sudo k3s ctr images import -
kubectl apply -f k8s/deployment.yaml -f k8s/service.yaml
kubectl rollout restart deployment cadencereader
