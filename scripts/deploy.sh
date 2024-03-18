#!/bin/bash

sudo docker network create shared-api-network || true
sudo docker compose up --build -d
