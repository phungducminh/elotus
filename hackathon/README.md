## Features
- Register
- Login
- Expired Token
- Upload file 
    - Check access
    - Check file size
    - Check file type is image
    - Store file in tmp/ folder

## Features not finished yet
- Load config from file or commandline
    - default expires token time: 1m
- Html form for upload file

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

## Notes
- A Postman collection and environment are also attached for better testing 
-- Elotus.postman_collection.json
-- Elotus.postman_environment.json

- The html form seem hard for testing since we also need to manage accessToken, 
which require login, and it's possible that I'm not available time for that yet

- DB schema
```
mysql> select * from users;
+----+----------+--------------------------------------------------------------+
| id | username | hashed_password                                              |
+----+----------+--------------------------------------------------------------+
|  1 | elotus   | $2a$10$5ReY3VPSEtUNb2Qu2/LSKOBpW1DBAdAKx9LZCbVB4BNhjp2qwp/G6 |
|  2 | elotus1  | $2a$10$cgzlXx8P0xLGtxj3iyAQ9uvbxo66NS7u55JS7yrcKl12xBnkPz1sO |
|  3 | elotus2  | $2a$10$CFEWYOt/dds5KRXRkpu1oul2uXP9Q6.z8XBriW5u2V9Su4dsCsolC |
+----+----------+--------------------------------------------------------------+

mysql> select * from files;
+----+---------+------------------------------------------+--------------+---------+
| id | user_id | filename                                 | content_type | size    |
+----+---------+------------------------------------------+--------------+---------+
|  1 |       1 | pexels-pixabay-358532.jpg                | image/jpeg   | 2594448 |
|  2 |       1 | pexels-pixabay-358532.jpg                | image/jpeg   | 2594448 |
|  3 |       1 | pexels-christian-heitz-285904-842711.jpg | image/jpeg   | 1847928 |
|  4 |       1 | pexels-souvenirpixels-417074.jpg         | image/jpeg   | 1275751 |
|  5 |       3 | pexels-eberhardgross-534164.jpg          | image/jpeg   | 4025828 |
+----+---------+------------------------------------------+--------------+---------+
```

-----------------------------------------------------------------
### Tools is used for development
- sqlc: generate type-safe code from SQL https://github.com/sqlc-dev/sqlc
    - only for development
- sql-migrate: SQL schema migration tool for Go. https://github.com/rubenv/sql-migrate
    - only for development (I've add create db script manually without sql-migrate)

### Enhanced features:

- add redis for storing user session
- add rate-limiter
- set up database: master and slave setup
- add distributed locking for username uniqueness
- add grpc gateway
- add AWS queue
- add circut breaker
- add retry
