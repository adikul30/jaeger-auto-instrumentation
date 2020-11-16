#!/bin/bash
set -x

cd serviceA
docker build -t adikul30/client-service:latest -f serviceADockerfile .
docker push adikul30/client-service:latest