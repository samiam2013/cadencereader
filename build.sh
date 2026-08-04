#!/bin/bash
set -e

HASH_FILE=".last-build-hash"

# Compute a hash of all .go files + go.mod/go.sum
CURRENT_HASH=$(find . -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \
  | sort \
  | xargs sha256sum \
  | sha256sum \
  | cut -d' ' -f1)

PREVIOUS_HASH=$(cat "$HASH_FILE" 2>/dev/null || echo "")

if [ "$CURRENT_HASH" != "$PREVIOUS_HASH" ]; then
  echo "Go source changed, rebuilding..."
  go build . || { echo preemptive go build failed; exit 1; }
  docker build -t cadencereader:latest .
  docker save cadencereader:latest | sudo k3s ctr images import -
  echo "$CURRENT_HASH" > "$HASH_FILE"
else
  echo "No changes to Go files, skipping build/rebuild."
fi

kubectl apply -f k8s/
kubectl rollout restart deployment cadencereader
