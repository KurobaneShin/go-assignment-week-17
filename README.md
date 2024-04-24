# go-assignment-week-17

# Server

This is a simple HTTP server built in Go that provides endpoints for handling in-memory data and fetching data from a MongoDB database.

## Features

- **In-memory Data Handling**: Provides endpoints for getting and setting key-value pairs in memory.
- **MongoDB Data Fetching**: Allows fetching data from a MongoDB collection based on specified criteria.

## Endpoints
### In-memory Data
GET /in-memory?key=<key>: Retrieve the value associated with the specified key.
POST /in-memory: Set a key-value pair in memory.

### MongoDB Data
POST /fetch-data: Fetch records from a MongoDB collection based on specified criteria.

## Usage
### Set a key-value pair
curl -X POST -d '{"key":"example", "value":"data"}' http://localhost:3334/in-memory

### Retrieve the value associated with a key
curl http://localhost:3334/in-memory?key=example

###Fetch data from MongoDB
curl -X POST -d '{"startDate":"2022-01-01", "endDate":"2022-01-31", "minCount":10, "maxCount":100}' http://localhost:3334/fetch-data

## Setup
make

### Prerequisites
- Go (1.13 or newer)
