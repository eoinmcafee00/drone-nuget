#!/bin/sh

scripts/build.sh

docker build -t drone/drone-nuget -f docker/Dockerfile .

docker run --rm -e PLUGIN_NUGET_APIKEY="oy2ey73rt5zc522ku2nqyr6iqf7e3tau3pxrdukyripq4a" \
  -e PLUGIN_NUGET_URI="https://api.nuget.org/v3/index.json" \
  -e PLUGIN_PACKAGE_LOCATION="test/CoreLibraries.1.0.0.nupkg" \
  -e DRONE_COMMIT_SHA=8f51ad7884c5eb69c11d260a31da7a745e6b78e2 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_BUILD_NUMBER=43 \
  -e DRONE_BUILD_STATUS=success \
  -w /drone/src \
  -v $(pwd):/drone/src \
  drone/drone-nuget