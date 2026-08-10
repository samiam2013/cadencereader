#!/bin/bash

kubectl exec -it example-db-1 -c postgres -- psql -U postgres
