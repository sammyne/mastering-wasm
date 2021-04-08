#!/bin/bash

echo "pre-deploy @localhost:8090"

workdir=$PWD/chapter01/code

docker run --rm \
  -v $workdir:/usr/share/nginx/html \
  --name html -p 8090:80 \
  nginx:1.19.6-alpine
