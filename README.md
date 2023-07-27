# Zero-Knowledge Proof (ZKP) Authentication Protocol

This repository implements a Zero-Knowledge Proof (ZKP) authentication protocol as a Proof-of-Concept application. The ZKP protocol is a viable alternative to password hashing in an authentication schema. The main goal of this project is to support one-factor authentication, which involves exact matching of a number (registration password) stored during registration and another number (login password) generated during the login process.

## Requirements

* golang v1.20
* protoc compiler (v23.3)
* docker (v20.10.21)
* docker-compose (v20.20.2)
* VSCode or any other suitable IDE

## Project Structure

The project is structured as follows:

```
zkp-authentication/
├── api
│   ├── v1
│   │   ├── err
│   │   │   └── error.go
│   │   └── proto
│   │       ├── zkp_auth_grpc.pb.go
│   │       ├── zkp_auth.pb.go
│   │       └── zkp_auth.proto
│   └── v2
│       ├── err
│       │   └── error.go
│       └── proto
│           ├── zkp_auth_grpc.pb.go
│           ├── zkp_auth.pb.go
│           └── zkp_auth.proto
├── cmd
│   └── cmd.go
├── deploy
│   └── local
│       ├── docker-compose.yml
│       ├── Dockerfile.client
│       └── Dockerfile.server
├── docs
│   ├── client.html
│   ├── index.html
│   ├── OVERVIEW.md
│   ├── server.html
│   └── zkp.html
├── go.mod
├── go.sum
├── internal
│   ├── client
│   │   └── client.go
│   ├── cpzkp
│   │   ├── cp_zkp.go
│   │   └── cp_zkp_test.go
│   ├── server
│   │   └── server.go
│   └── tests
│       ├── client_test.go
│       └── server_test.go
├── lib
│   ├── config
│   │   └── config.go
│   └── util
│       └── util.go
├── LICENSE
├── main.go
├── Makefile
└── README.md
```

## Usage

1. Clone the repository:

```
git clone https://github.com/username/zkp-authentication.git
```

2. Change into the project directory:

```
cd zkp-authentication
```

3. Build the binaries:

```
go build .
```

4. Start the ZKP authentication server:

```
./zkp-authentication --server
```

5. Use the CLI to register a new user:

```
./zkp-authentication register -u <username> -p <password>
```

6. Use the CLI to log in with the registered user:

```
./zkp-authentication login -u <username> -p <password>
```

Alterntively, if you don't wish to build the binaries, then


4. Start the ZKP authentication server:

```
go run main.go --server
```

5. Use the CLI to register a new user:

```
go run main.go register -u <username> -p <password>
```

6. Use the CLI to log in with the registered user:

```
go run main.go login -u <username> -p <password>
```



## Testing

### Unit Tests

To run the unit tests, use the following command:

```
make test
```

### Run with Docker

To run using Docker, ensure that Docker is installed on your machine and follow these steps:

1. Build the Docker images and containers:

```

docker compose up
```

2. Enter the Docker `local-zkp-auth-client` container:

```
docker exec -it local-zkp-auth-client sh
```

Repeat steps 4, 5, and 6 under the **Usage** section

3. Stop and remove the Docker containers:

```
docker-compose down
```


## License

This project is licensed under the [MIT License](LICENSE).
