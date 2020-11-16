#!/bin/bash
set -x

cd serviceB
docker build -t adikul30/server-service:latest -f serviceBDockerfile .
docker push adikul30/server-service:latest