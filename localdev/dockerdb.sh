#!/bin/bash

docker run -d --name devdb \
  -e POSTGRES_PASSWORD=dev \
  -e POSTGRES_DB=cadencereader \
  -p 5432:5432 \
  -v devdb-data:/var/lib/postgresql/data \
  postgres:16

migrate -path ../migrations -database "postgres://postgres:dev@localhost:5432/cadencereader?sslmode=disable" up
