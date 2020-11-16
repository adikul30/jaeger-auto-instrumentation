#!/bin/bash
set -x

cd proxy
docker build -t adikul30/sidecar-proxy:latest -f proxyDockerfile .
docker push adikul30/sidecar-proxy:latest