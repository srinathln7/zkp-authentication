# package `client`

The provided code implements a gRPC client for the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) authentication protocol. The client is responsible for registering users and performing logins with the server using ZKP-based authentication.

1. **Imports and Package:**
   - The code imports necessary packages and defines the `client` package.
   - It includes gRPC-related packages, color formatting, error handling, CP-ZKP package, and utility functions.

2. **Type Definitions:**
   - `RegRes` and `LogInRes` are structs to store registration and login responses.

3. **SetupGRPCClient Function:**
   - `SetupGRPCClient` sets up the gRPC client and returns the `AuthClient`.
   - It loads the server address from the `.env` file.
   - The gRPC connection is established with insecure credentials.
   - The gRPC client is created with the established connection.

4. **Register Function:**
   - `Register` handles user registration with the server using ZKP.
   - It generates the CP-ZKP system parameters (`cpzkpParams`) by creating a new `CPZKP` instance.
   - The user's password is uniquely converted to a big integer `x` using the `getSecretValue` function.
   - A new prover (client) is created based on `x` (secret value), and it calculates `y1` and `y2` values.
   - The client sends the registration request to the server with the calculated `y1` and `y2`.
   - If successful, it returns a registration response message.

5. **LogIn Function:**
   - `LogIn` performs user login with the server using ZKP.
   - It generates the CP-ZKP system parameters (`cpzkpParams`) by creating a new `CPZKP` instance.
   - The user's password is uniquely converted to a big integer `x` (secret value) using the `getSecretValue` function.
   - A new prover (client) is created based on `x`, and it calculates commitment values `r1` and `r2`.
   - The client sends the authentication challenge request to the server with `r1` and `r2`.
   - The server responds with an authentication challenge, including `authID` and `c`.
   - The client calculates the response `s` using the received `c` and the prover's secret value `x`.
   - The client verifies the authentication response with the server by sending `authID` and `s`.
   - If successful, it returns a login response with a session ID.

6. **getSecretValue Function:**
   - `getSecretValue` converts a password string to a unique big integer using the utility library function `StringToUniqueBigInt`. For more info on the functions in the utility 
   library refer [here](https://github.com/srinathLN7/zkp-authentication/tree/main/lib/util).

The CP-ZKP client code provides a gRPC-based authentication client that allows users to register and login securely using the Chaum-Pedersen Zero-Knowledge Proof protocol. The client generates and sends ZKP-based proof commitments and responses to the server for authentication. It also includes error handling for invalid requests and responses. The client works with the CP-ZKP server to securely perform user registration and login operations.
