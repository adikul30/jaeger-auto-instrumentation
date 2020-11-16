#!/bin/bash
set -x

cd init
docker build -t adikul30/init-container:latest -f initDockerfile .
docker push adikul30/init-container:latest