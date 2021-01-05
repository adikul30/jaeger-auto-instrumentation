#!/bin/bash
set -x

./utils/buildinit.sh
./utils/buildproxy.sh
./utils/buildclient.sh
./utils/buildserver.sh
