# Pacakge `server`

The provided code defines a gRPC server that implements the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol for user registration and authentication. The server stores user information and performs ZKP-based authentication.

1. **Imports and Package:**
   - The code imports necessary packages and defines the `server` package.
   - It includes gRPC-related packages, error handling, UUID generation, and the CP-ZKP package.

2. **Type Definitions:**
   - `CPZKP` interface represents the methods required for initializing CP-ZKP parameters.
   - `Config` struct holds the CP-ZKP configuration.
   - `RegParams` and `AuthParams` are structs used to store registration and authentication parameters for users.

3. **`grpcServer` Struct:**
   - `grpcServer` is the main struct representing the CP-ZKP server.
   - It includes fields for the user registration directory (`RegDir`) and authentication directory (`AuthDir`).
   - `*Config` holds the CP-ZKP configuration.

4. **RunServer Function:**
   - `RunServer` function is the entry point of the server.
   - It loads the configuration from the `.env` file using `godotenv`.
   - The server is created, and the gRPC server is started with the specified address and port.

5. **`newgrpcServer` Function:**
   - `newgrpcServer` initializes the `grpcServer` with the CP-ZKP system parameters and empty user directories.

6. **NewGRPCServer Function:**
   - `NewGRPCServer` creates a new gRPC server, registers the service, and returns the server.

7. **Register Function:**
   - `Register` handles user registration on the server.
   - It checks if the user is already registered (`RegDir`).
   - If not, it parses and stores the provided `y1` and `y2` values for every unique user in the registration directory.
   - If the user is already registered, it returns an error indicating an invalid registration.

8. **CreateAuthenticationChallenge Function:**
   - `CreateAuthenticationChallenge` handles the authentication challenge generation for registered users.
   - It checks if the user is registered.
   - If the user is registered, it creates a verifier, generates a challenge (`c`), and stores it againt the unique `auth_id` (UUID) in authentication directory.
   - The `auth_id`, along with `c`, is returned in the response.

9. **VerifyAuthentication Function:**
   - `VerifyAuthentication` verifies the user's response to the authentication challenge.
   - It checks the validity of the provided `auth_id`.
   - If the `auth_id` is valid, it retrieves the user's information and the stored challenge (`c`) from `AuthDir`.
   - The user's (`y1`, `y2`) and (`r1`,`r2`) values are also retrieved from `RegDir` and `AuthDir` respectively.
   - The user's response `S` is parsed into a big integer.
   - A verifier is created, and the proof is verified using `VerifyProof`.
   - If the proof is valid, a session ID (UUID) is generated and returned in the response. Otherwise, a 401 authentication error is thrown with details.


The above server code provides a gRPC-based authentication service using the Chaum-Pedersen Zero-Knowledge Proof protocol. It allows users to register their `y1` and `y2` values and subsequently authenticate using the ZKP protocol. The server verifies the correctness of the authentication challenge and generates a session ID for authenticated users. The server simulates user registration and authentication using in-memory non-persistent storage.
