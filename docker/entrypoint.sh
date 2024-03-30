#!/bin/sh

envsubst < ./etc/market-api.yaml.tmpl > ./etc/market-api.yaml
./market-server
