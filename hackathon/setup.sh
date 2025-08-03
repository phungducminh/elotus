#!/usr/bin/env bash

docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=elotus --mount type=volume,src=mysql,dst=/var/lib/mysql mysql:8.0.20
docker exec -it mysql mysql -h 127.0.0.1 -u root -p'elotus' -e 'create database elotus;'
