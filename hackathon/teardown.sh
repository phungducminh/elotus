#!/usr/bin/env bash

docker stop mysql
docker rm mysql
docker volume rm mysql
