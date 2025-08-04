## How to run the service
- Requirements: Docker is available
### Steps (In MacOS Intel)
- Make sure current directory is /hackathon
- Run setup.sh: it should run a mysql docker container and create "elotus" db
```
bash ./setup.sh
```
- Run command (I've tried to put this command into setup.sh file, but failed)
```
docker exec -it mysql /bin/sh -c "mysql -h 127.0.0.1 -u root -p'elotus' < /dbscript.sql"
```
- Run the Golang service (might require to run "go mod tidy")
```
go run main.go
```
- Sample test cases
```
curl --location '127.0.0.1:8080/api/auth/register' \
--header 'Content-Type: application/json' \
--data '{
    "username": "elotus",
    "password": "elotus"
}'

curl --location '127.0.0.1:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "elotus",
    "password": "elotus"
}'

curl --location '127.0.0.1:8080/api/file/upload' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQyNzM1MTUsImlhdCI6MTc1NDI3MzQ1NSwiaXNzIjoiZWxvdHVzIiwic3ViIjoiMSJ9.NoHOH4YA1nn3JR5HMkmexcfjtDNUvHabyAyw2TkSk4g' \
--form 'data=@"<PLACEHOLDER>"'
```
- To reset database, please run:
```
bash ./teardown.sh
```

-----------------------------------------------------------------
### Enhanced features:

- add redis for storing user session
- add rate-limiter
- set up database: master and slave setup
- add distributed locking for username uniqueness
- add grpc gateway
- add AWS queue
- add circut breaker
- add retry
