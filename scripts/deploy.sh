#!/bin/bash

sudo docker network create shared-api-network || true
sudo docker compose up --build -d --remove-orphans

RETRY_INTERVAL=1
MAX_RETRIES=15
RETRIES=0

until (( RETRIES == MAX_RETRIES )) || curl -sSf http://localhost/api/heartbeat; do
    (( RETRIES++ ))
    sleep $RETRY_INTERVAL
done

(( RETRIES >= MAX_RETRIES ))
