# Required features
## create an simple form (with default server url) and upload file via that input
## Create database from the script

# Enhanced features:

## add redis
### fetch user and store user information
### expiration time

## add rate-limiter
### for register and login endpoints
### buffered chan for upload file?

## set up database: master and slave setup
## TODO: add distributed locking for username uniqueness
## TODO: replace logger
## TODO: add grpc gateway
## TODO: add AWS queue
## TODO: add docker
## TODO: add circut breaker
## TODO: add retry


# Estimation
## 100M users -> 20M DAU, 
### Login 2 times/d -> 40 M req/d -> 400 reqs/s, peak 1000 req/s, if we store user in cache 2d -> peak less than 1000
### Register: 1M reqs/d -> 10 reqs/s 
### max_connections is 1000 is good enough (not that we didn't store login time, each time user login, otherwise it will be 1000 reqs/s to db for just storing accessToken)

# File upload
+ maxMemory: ParseMultipartForm(maxMemory) ??? what is the meaning of maxMemory?
+ what is the file part that can't be stored in memory?
maxFileMemoryBytes: limit for only file content
file metadata is always stored to memory
