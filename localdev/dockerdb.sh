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

set -a
source .env
set +a

db_container_name="$DB_NAME-container"
db_volume_name="$DB_NAME-data"

docker stop $db_container_name
docker rm $db_container_name
docker volume rm $db_volume_name

docker run -d --name $db_container_name \
  -e POSTGRES_PASSWORD=$DB_PASS \
  -e POSTGRES_DB=$DB_NAME \
  -p 5432:5432 \
  -v $db_volume_name:/var/lib/postgresql/data \
  postgres:16

until docker exec $db_container_name psql -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; do
    sleep 0.1
done
echo "Postgres is ready!"

migrate -path ../database/migrations -database "postgres://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME?sslmode=disable" up
