# Redis Clone

Redis Clone is an educational project to learn and understand the inner workings of Redis, a popular in-memory data structure store. This project is implemented in Go (Golang) and covers various aspects of Redis, including in-memory storage, data persistence, and advanced features like pub/sub and transactions.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Introduction
Redis Clone aims to mimic the basic functionalities of Redis to provide a learning platform for developers interested in understanding distributed systems, data structures, and high-performance computing.

## Features
- In-memory key-value store
- Data persistence using RDB and AOF
- Support for lists, sets, and hash maps
- Publish/Subscribe messaging
- Basic transaction support
- Lua scripting execution
- Master-slave replication (planned)
- Sharding (planned)

## Installation
To get started with Redis Clone, follow these steps:

1. **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/redis-clone.git
    cd redis-clone
    ```

2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

3. **Build the project**:
    ```bash
    go build -o redis-server ./cmd/redis-server
    ```

## Usage
Run the Redis Clone server:

```bash
./redis-server
```
You can then interact with the server using any Redis client.

## License
This project is licensed under the MIT License.

## Acknowledgements
- [Redis](https://redis.io/documentation) for the inspiration and original implementation.
- [Golang](https://golang.org/) for the programming language.
