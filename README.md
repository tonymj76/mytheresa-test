# Mytheresa Promotion Assessment


## Introduction
This REST API endpoint transforms a list of products into a format suitable for display. It applies discounts based on 
specific criteria being met.

This is the backend repository. It's written in Go.
The Architecture use is Domain-Driven Design (DDD) because DDD emphasizes understanding the domain and building the software around it, ensuring the system solves the real business problems effectively.

### Technology Used
- Go (go-gin)
- Docker and Docker-compose
- testify for testing
- Postgres
- Ent ORM library
- make file


## SetUp

1. After clone the repo with the following command
    ```bash
    git clone git@github.com:tonymj76/mytheresa-test.git
    ```
2. change directory to the folder `mytheresa-test`
3. Rename the following file from `.env.example` to `.env` Eg `mv .env.example .env` in your terminal


To run the application please have docker and docker-compose install. If you also have make install you can use the command in the make file [here](makefile)

- run `docker-compose up --build -d`
- the server is running on port 9191
- the first endpoint to send request is `http://localhost:9191/api/products`
```
GET /products                                       // Read all products and apply discounts
GET /products?category=boots                        // Read product that belong in boots category and apply discount if the criteria are met
GET /products?priceLessThan=89000                  // Read product with priceLessThan=89000 which will get price <= 89000
GET /products?category=boots&priceLessThan=89000    // category filtering takes precedence here. which will ignore priceLessThan=89000 
```

## To run Test
 ```
 go test ./handlers -run=Handler -v
 ```




if you have make installed you can run the following commands for shortcuts 
use `make run` to run in container and `make log` to show logs or `make test` to run test.

for further reading on [entgo](https://entgo.io/docs/getting-started/)
