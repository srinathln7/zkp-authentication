# Zero-Knowledge Proof (ZKP) Authentication Protocol

This repository implements a Zero-Knowledge Proof (ZKP) authentication protocol as a Proof-of-Concept application. The ZKP protocol is a viable alternative to password hashing in an authentication schema. The main goal of this project is to support one-factor authentication, which involves exact matching of a number (registration password) stored during registration and another number (login password) generated during the login process.

## ZKP Protocol Overview

The ZKP protocol used in this project is based on the Chaum–Pedersen Protocol described in the book ["Cryptography: An Introduction (3rd Edition) Nigel Smart"](https://www.cs.umd.edu/~waa/414-F11/IntroToCrypto.pdf)

### Registration Process

In the registration process, the prover (client) has a secret password `x` (a number) and wants to register it with the verifier (server). The prover calculates `y1` and `y2` using public parameters `g` and `h`, along with the secret `x`, and sends `y1` and `y2` to the verifier.

### Login Process

The login process follows the ZKP Protocol, where the prover is the authenticating party, and the verifier is the server running the authentication check.

## Project Goals

The primary goal of this project is to design and write the code that implements the ZKP Protocol outlined above. The solution should be implemented as a server and client using the gRPC protocol according to the provided interface described in the `protobuf` schema. 

## Additional Features

- Unit tests, where appropriate
- Functional tests of the ZKP Protocol
- A setup to run the Client and the Server
- Performance optimizations
- Coverage of test cases (not code coverage)
- Code soundness
- Code organization
- Code quality
- Well-documented code
- Each instance runs in a separated Docker container with a Docker Compose setup
- Code to deploy the two containers in AWS (client on one machine and server on another machine)
- Implementation of two flavors: One with exponentiations (as described in the book) and one using Elliptic Curve cryptography (similar to the ZKP implementation in Rust)
- Allow using "BigInt" numbers for increased precision

