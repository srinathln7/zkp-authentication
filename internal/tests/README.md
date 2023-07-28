# Package `test`

The `test` package consists of two file `client_test.go` and `server_test.go` file to test our grpc client and server implementation.

## `client_test.go`:

1. **SetupGRPCClient Function:**
   - This function sets up the gRPC client for testing purposes.
   - It generates a free available port (`listener`) for the gRPC server.
   - The client and server communicate via an insecure connection (`insecure.NewCredentials()`).
   - CP-ZKP parameters (`cpzkpParams`) are generated for the server configuration (`cfg`).
   - If the optional function `fn` is provided, it initializes the server with the provided configuration.
   - The gRPC server is started in a separate goroutine using `go func()`.
   - The function returns the gRPC client, server configuration, and a `teardown` function to stop the server and close connections.

2. **testClientRegisterUserSuccess Function:**
   - This function tests the successful user registration on the server.
   - It generates CP-ZKP system parameters (`cpzkpParams`) and a correct secret value `x` for the prover (client).
   - The prover generates `y1` and `y2` values based on the secret value and CP-ZKP system parameters.
   - The client sends a registration request to the server with the username `srinath` and the generated `y1` and `y2`.
   - The function verifies that the server responds with an empty registration response (indicating successful registration).

3. **testClientRegisterUserFail Function:**
   - This function tests the failure scenario for user registration on the server.
   - It generates CP-ZKP system parameters (`cpzkpParams`) and a correct secret value `x` for the prover (client).
   - The prover generates `y1` and `y2` values based on the secret value and CP-ZKP system parameters.
   - The client attempts to register the user `srinath` again with the same `y1` and `y2`.
   - The function expects a non-nil error from the server, indicating that registration fails incase of duplicate registration.

4. **testClientVerifyProofSuccess Function:**
   - This function tests the successful generation and verification of a proof by the client.
   - It generates CP-ZKP system parameters (`cpzkpParams`) and a correct secret value `x` for the prover (client).
   - The prover generates `r1` and `r2` values based on the secret value and CP-ZKP system parameters.
   - The client sends an authentication challenge request to the server with the username `srinath` and the generated `r1` and `r2`.
   - The server responds with an authentication challenge, including an `authID` and `c`.
   - The prover calculates the response `s` using the received `c` and the secret value `x`.
   - The client sends the authentication response (`s`) to the server for verification.
   - The function verifies that the server accepts the proof, indicating successful authentication.

5. **testClientVerifyProofFail Function:**
   - This function tests the failure scenario for proof verification by the client.
   - It generates CP-ZKP system parameters (`cpzkpParams`) and an incorrect secret value `x` for the prover (client).
   - The prover generates `r1` and `r2` values based on the incorrect secret value and CP-ZKP system parameters.
   - The client sends an authentication challenge request to the server with the username `srinath` and the generated `r1` and `r2`.
   - The server responds with an authentication challenge, including an `authID` and `c`.
   - The prover calculates an incorrect response `s` due to the incorrect secret value.
   - The client sends the authentication response (`s`) to the server for verification.
   - The function expects an error from the server that matches the expected `grpc_err.ErrInvalidChallengeResponse` defined in the `api/v2/err/error.go` file.

## `server_test.go`:

1. **TestMain Function:**
   - The `TestMain` function is an entry point for running tests in this package.
   - It calls `m.Run()` to execute the tests in the package.

2. **TestGRPCServer Function:**
   - The `TestGRPCServer` function is the main test function for testing the gRPC server.
   - It sets up the gRPC client and server using `SetupGRPCClient`.
   - The server configuration and teardown function are also returned from `SetupGRPCClient`.
   - Subtests are used to group individual test cases, providing better isolation, readability, and reporting.
   - The following test cases are run as subtests:
     - `register user successfully`: Tests the successful user registration on the server.
     - `verification proof successfully`: Tests the successful proof generation and verification by the client.
     - `verification proof failure`: Tests the failure scenario for proof verification by the client.
     - `register user failure`: Tests the failure scenario for user registration on the server.
   - The server is gracefully shutdown and connections are closed after finishing all the test cases using `teardown`.


The test files (`client_test.go` and `server_test.go`) contain comprehensive test functions for the CP-ZKP client and server implementations. These tests ensure that the client and server work correctly in various scenarios, including successful and failed user registration, successful proof generation and verification, and failed proof verification. The use of subtests allows for better organization and isolation of the test cases, making it easier to pinpoint and resolve potential issues.
