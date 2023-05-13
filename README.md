# Batas Platform API

## Installation

Requirements:

* [Go 1.16](https://go.dev/dl/)
* [Redis server](https://redis.io/download/)

Installation Steps:

* Run `go install github.com/swaggo/swag/cmd/swag@latest`
* Run `go mod tidy`
* Duplicate `.env.example` to `.env` and fill in mainly MySQL and Redis configuration

That's it.

## Running the script

To run the script simply type `go run main.go`

## Generating Swagger Doc

Make sure you're in project root folder and run `swag init`. Swagger documentation will be generated on `docs` folder. Run `go run main.go` and open http://127.0.0.1:8000/swagger/index.html to preview the Swagger documentation.