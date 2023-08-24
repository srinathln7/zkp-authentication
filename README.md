# Zero-Knowledge Proof (ZKP) Authentication Protocol

This repository implements a Zero-Knowledge Proof (ZKP) authentication protocol as a Proof-of-Concept application. The ZKP protocol is a viable alternative to password hashing in an authentication schema. The main goal of this project is to support one-factor authentication, which involves exact matching of a number (registration password) stored during registration and another number (login password) generated during the login process. 

Click [here](https://www.youtube.com/watch?v=-ueOQ7y35Ms) to watch the demo presentation. 

## Requirements

* golang (v1.20)
* protoc compiler (v23.3)
* docker (v20.10.21)
* docker-compose (v20.20.2)
* VSCode or any other suitable IDE

## Project Structure

Refer [here](https://github.com/srinathLN7/zkp-authentication/blob/main/OVERVIEW.md) for the complete overview of the protocol and the project structure.

## Approach

To achieve the implementation of the Zero-Knowledge Proof (ZKP) authentication protocol, the following approach is taken:

### Building Relevant gRPC Server and Client APIs:
   - The gRPC server and client APIs are constructed based on the given `protobuf` schema using the `protoc` compiler. 

### Implementing Chaum-Pedersen Zero Knowledge Proof Protocol:
   - The Chaum-Pedersen Zero Knowledge Proof Protocol is implemented and tested in isolation. Please note that in order to support Big integers, the variables `r1`, `r2`, `c`, and `s` in the 
Zero-Knowledge Proof (ZKP) authentication protocol have been changed from `int64` to `string`. This change was necessary because `int64` data type has a fixed range of representable numbers (`-2^63` to `2^63-1`), and it may not be able to handle large integers that are required for cryptographic operations. By using the `string` data type, the ZKP protocol can now accommodate big integers without any limitation on their size. This ensures that the protocol remains secure and accurate even when dealing with large cryptographic values. With this update, the ZKP authentication protocol is better equipped to handle the complexities of cryptographic operations and provide a more reliable and secure user authentication process.For implementaion details of the protocol, refer [here](https://github.com/srinathLN7/zkp-authentication/tree/main/internal/cpzkp).

### Building the gRPC Server:
   - The gRPC server is built using the `protoc` generated `zkp_auth_grpc.pb.go` and `zkp_auth.pb.go` files. For detailed information, check [here](https://github.com/srinathLN7/zkp-authentication/tree/main/internal/server).

### Building the gRPC Client:
   - The gRPC client is developed using the `protoc` generated `zkp_auth_grpc.pb.go` and `zkp_auth.pb.go` files. For more information, click [here](https://github.com/srinathLN7/zkp-authentication/tree/main/internal/client).

### Testing the Server and Client:
   - Comprehensive testing is performed on the built gRPC server and client to ensure their functionality. For more details, see [here](https://github.com/srinathLN7/zkp-authentication/tree/main/internal/tests).

### Command Line Interface (CLI) Application:
   - A command-line interface (CLI) application is developed to provide a user-friendly interface for the ZKP authentication protocol. For more details about the CLI, check [here](https://github.com/srinathLN7/zkp-authentication/tree/main/cmd).
   
### Main Package Entry Point:
   - The `main` package serves as the entry point for the entire program, orchestrating the ZKP authentication protocol's execution. Refer [here](https://github.com/srinathLN7/zkp-authentication/blob/main/MAIN.md) for more details on the same.
   
For a more in-depth understanding of each step, please refer to the relevant documentation provided in the links.

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
./zkp-auth register -u <username> -p <password>
```

6. Use the CLI to log in with the registered user:

```
./zkp-auth login -u <username> -p <password>
```

Alternatively, if you don't wish to build the binaries, then

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

To run all the test files in this project, run the following command in your local development terminal:

```
make test
```

Upon running this command, you should see all the test cases passing, ensuring the proper functioning of all components within our project. Successful test results indicate that the application is operating as expected and meeting the desired requirements. 


## Run with Docker

To run using Docker, ensure that Docker is installed on your machine and follow these steps:

1. Build the Docker images and containers:

```

cd deploy/local

docker compose up

```

If you encounter issues with building the containers due to IP address overlap, it is likely caused by conflicting IP addresses in the network. To resolve this, you can change the subnet address used for the containers to ensure uniqueness. By selecting a different subnet address, you can avoid conflicts and successfully build the containers.

2. Enter the Docker `local-zkp-auth-client` container:

```
docker exec -it local-zkp-auth-client sh
```

Repeat steps 5, and 6 under the **Usage** section to test for various usernames and password. 

3. Stop and remove the Docker containers:

```
docker compose down
```



## API Documentation

For the API documentation, refer to the [docs](https://github.com/srinathLN7/zkp-authentication/tree/main/docs) directory containing individual API documentation about the gRPC server and client APIs.

Alternatively, if you wish to build your own docs, run:

```
godoc 
```

and navigate to http://localhost:6060/pkg/github.com/srinathLN7/zkp_auth/internal/?m=all in your browser. You will find the links to all three packages: server, client, and cpzkp (Chaum-Pedersen ZKP).


## Improvements

To enhance the Zero-Knowledge Proof (ZKP) authentication protocol, the following improvements can be implemented:

* Implement Mutual TLS Authentication:
  - Introduce Mutual TLS-based authentication between the gRPC server and client to establish a secure and trusted communication channel. This ensures that both parties can verify each other's identities and encrypt the data exchanged during communication.

* Deployment Scripts for AWS:
  - Develop deployment scripts to automate the process of deploying the gRPC server and client containers to the Amazon Web Services (AWS) platform. This streamlines the deployment process and facilitates scalability and reliability.

* Chaum-Pedersen ZKP using Elliptic Curve Cryptography (ECC):
  - Enhance the Chaum-Pedersen Zero Knowledge Proof Protocol by implementing it using Elliptic Curve Cryptography (ECC). 

* Integration of SQL/NoSQL Database:
  - Integrate a SQL or NoSQL database into the ZKP authentication protocol to enable persistent data storage. Storing user data in a database ensures that user information is retained across server restarts and provides better support for user management and authentication.

* Enhanced Character Support:
  - Currently, the transformation of the user password into a secret value `x` is performed using the `StringToUniqueBigInt` function in the `util` library. This function relies on ASCII characters and uses base 256 to encode the password string. However, to accommodate passwords containing characters beyond the ASCII range, we can extend the character support by utilizing a wider range of characters for encoding.



## Self Review: Post Submission Comments   

This branch contains some suggestions to improve the codebase in `main` branch

**Done**
* Updated the `SERVER_ADDRESS` environment variable for both the `zkp-auth-server` and `zkp-auth-client` containers to use the service name `zkp-auth-server` as the hostname. Docker Compose's built-in DNS resolution will automatically resolve this hostname to the IP address of the corresponding container.
* Remove unnecessary declaration of channel `c` in `cmd.go` file inside `RootCmd` function. 
* `log.Fatal`, `log.Fatalf` internally calls `os.Exit(1)`. Hence all `os.Exit(1)` statements can be removed after `log.Fatal` and `log.Fatalf` statements.

**Tobe Done**
* Introduce custom grpc error `ErrUSerNotFound` with status code `404` in the `error.go` file to cover the case where the client invokes `login` before registering a user.
* Combine the two test files `server_test.go` and `client_test.go` inside the `internal/test` directory into one single file `grpc_test.go`. Consider renaming
the function `setupGRPCClient` to `setupGRPCTest`.

## License

This project is licensed under the [MIT License](LICENSE).
