# singkatin-revamp
Revamped [Last Projects URL Shortener API's](https://github.com/PickHD/singkatin-api) using Golang.

## Architect Overview :
![ARCH](https://raw.github.com/PickHD/singkatin-revamp/master/arch_singkatin_revamp.png)

## Whats Changed? :
1. From Monolith to Microservices
2. Scalable & Maintainable Codes
3. Become Asynchronous Operation using Messaging Queues (RabbitMQ)
4. From SQL to NoSQL
5. Using RPC for service interactions

## Main Features : 
1. Register
2. Login
3. Reset Password
4. User Profiles
5. User Dashboard (can analyze how much visitor / users click the short links)
6. Shortener Link Redirect

## Tech Used :
1. Golang _(Every services using different framework due experimenting performances.)_
2. MongoDB
3. Redis
4. RabbitMQ
5. GRPC
6. Docker
7. Jaeger
8. MinIO Storage

## Prerequisites : 
1. Make sure Docker & Docker Compose already installed on your machine
2. Rename `example.env` to `.env` on folder `./cmd/v1` every services
3. Make sure to uncheck comment & fill your **SMTP configuration** on auth env

## Setup :
1. To build all services, run command : 
    ```
    make build
    ```

2. You can build & run all services in background using command : 
    ``` 
    make run
    ```
3. If you want to stop all services then run :
    ```
    make stop
    ```
4. Last if want to stop & remove entire services then run :
    ```
    make remove
    ```