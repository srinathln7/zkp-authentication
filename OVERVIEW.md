# ZKP Protocol Overview

The ZKP protocol used in this project is based on the Chaum–Pedersen Protocol described in the book ["Cryptography: An Introduction (3rd Edition) Nigel Smart"](https://www.cs.umd.edu/~waa/414-F11/IntroToCrypto.pdf)

## Registration Process

In the registration process, the prover (client) has a secret password `x` (a number) and wants to register it with the verifier (server). The prover calculates `y1` and `y2` using public parameters `g` and `h`, along with the secret `x`, and sends `y1` and `y2` to the verifier.

## Login Process

The login process follows the ZKP Protocol, where the prover is the authenticating party, and the verifier is the server running the authentication check.

## Project Goal

The primary goal of this project is to design and write the code that implements the ZKP Protocol outlined above. The solution should be implemented as a server and client using the gRPC protocol according to the provided interface described in the [`protobuf`](https://github.com/srinathLN7/zkp-authentication/blob/main/api/v2/proto/zkp_auth.proto) schema. 


## Project Overview

```
.
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
│       │   ├── error.go
│       │   └── README.md
│       └── proto
│           ├── zkp_auth_grpc.pb.go
│           ├── zkp_auth.pb.go
│           └── zkp_auth.proto
├── cmd
│   ├── cmd.go
│   └── README.md
├── deploy
│   └── local
│       ├── docker-compose.yml
│       ├── Dockerfile.client
│       └── Dockerfile.server
├── docs
│   ├── client.html
│   ├── index.html
│   ├── server.html
│   └── zkp.html
├── go.mod
├── go.sum
├── internal
│   ├── client
│   │   ├── client.go
│   │   └── README.md
│   ├── cpzkp
│   │   ├── cp_zkp.go
│   │   ├── cp_zkp_test.go
│   │   └── README.md
│   ├── server
│   │   ├── README.md
│   │   └── server.go
│   └── tests
│       ├── client_test.go
│       ├── README.md
│       └── server_test.go
├── lib
│   ├── config
│   │   └── config.go
│   └── util
│       ├── README.md
│       └── util.go
├── LICENSE
├── main.go
├── MAIN.md
├── Makefile
├── OVERVIEW.md
└── README.md

19 directories, 39 files

```
