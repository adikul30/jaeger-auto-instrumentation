#!/bin/bash
set -x

cd client
docker build -t adikul30/client-service:latest -f clientDockerfile .
docker push adikul30/client-service:latest