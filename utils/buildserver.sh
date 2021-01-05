#!/bin/bash
set -x

cd server
docker build -t adikul30/server-service:latest -f serverDockerfile .
docker push adikul30/server-service:latest