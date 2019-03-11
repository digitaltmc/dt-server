#!/bin/bash

set -e

function _dbshell() {
  docker-compose exec shell mongo --host db
}

[ -z "$1" ] && echo "Specify some params." && exit 1

case "$1" in
  up|down)
    docker-compose $1
    ;;
  *)
    _$1
esac

