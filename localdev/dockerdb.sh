#!/bin/bash

#Open Docker, only if is not running
if (! docker stats --no-stream ); then
  # On Mac OS this would be the terminal command to launch Docker
  open /Applications/Docker.app
 #Wait until Docker daemon is running and has completed initialisation
while (! docker stats --no-stream ); do
  # Docker takes a few seconds to initialize
  echo "Waiting for Docker to launch..."
  sleep 1
done
fi

docker run -d --name devdb \
  -e POSTGRES_PASSWORD=dev \
  -e POSTGRES_DB=cadencereader \
  -p 5432:5432 \
  -v devdb-data:/var/lib/postgresql/data \
  postgres:16

migrate -path ../migrations -database "postgres://postgres:dev@localhost:5432/cadencereader?sslmode=disable" up
