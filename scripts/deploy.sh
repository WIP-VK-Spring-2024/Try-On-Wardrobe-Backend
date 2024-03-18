#!/bin/bach

docker network create shared-api-network || true
docker compose up --build -d
