#!/usr/bin/env bash

docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=elotus --mount type=volume,src=mysql,dst=/var/lib/mysql mysql:8.0.20
docker cp ./dbscript.sql mysql:/

echo "Waiting for MySQL to start..."
until docker exec -it mysql mysqladmin ping -h localhost -u root -p'elotus' --silent; do
  sleep 1
done
echo "MySQL is ready, please run database dbscript.sql!"
