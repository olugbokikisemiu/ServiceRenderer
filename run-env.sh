#!/bin/sh

set -e

CALLBACK_HOST="localhost"

PATCH_YML=""

unameOut="$(uname -s)"
case "${unameOut}" in
    Darwin*)
      CALLBACK_HOST="docker.for.mac.localhost"
    ;;
esac

docker-compose \
    -f docker-compose.yml \
    rm -f --stop 1>&2 || true

docker-compose \
    -f docker-compose.yml \
    up -d --force-recreate 1>&2