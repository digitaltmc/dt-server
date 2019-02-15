#!/bin/bash

# docker pull mongo

function _default () {
docker volume create dt_data
docker run -d --rm --name dt_data \
  -e MONGO_INITDB_DATABASE=dt \
  -p 27017:27017 \
  mongo
}

# TODO: add auth to mongo
#  -e MONGO_INITDB_ROOT_USERNAME=admin \
#  -e MONGO_INITDB_ROOT_PASSWORD=password \


#---------- Main

if [[ -z "$1" ]]; then
  _default
fi

while getopts ":f:" opt; do
  case $opt in
    f)
      echo "-f was triggered with $OPTARG!" >&2
      _$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done
