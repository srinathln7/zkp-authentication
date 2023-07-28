# Package `cmd` 

The `cmd` package defines a CLI tool that enables users to register and login using Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) based authentication. The CLI provides two subcommands, `register` and `login`, which are used for user registration and login, respectively. The program sets up a gRPC client for communication with the server and prints the results of registration and login requests in JSON format with colored output. The CLI runs until a termination signal (CTRL+C) is received, and it exits gracefully.

1. **SetupFlags Function:**
   - This function is used to set up command-line flags for the CLI.
   - It defines two flags: `user` and `password`, which can be used as options for the `register` and `login` subcommands.
   - The flags are associated with the root command (`RootCmd`) and added to it.
   - Two subcommands, `registerCmd` and `loginCmd`, are also added to the root command.

2. **RootCmd:**
   - `RootCmd` is the main Cobra command for the CLI tool.
   - It has a `Run` function that prints a welcome message and instructions for using the CLI.
   - The instructions provide details on how to use the `register` and `login` subcommands.
   - The program waits for a termination signal (CTRL+C) to exit gracefully.

3. **registerCmd:**
   - `registerCmd` is a subcommand that represents the `register` functionality of the CLI.
   - When invoked, it sets up a gRPC client (`grpcClient`) for communication with the server.
   - The `client.SetupGRPCClient()` function is used to set up the gRPC client.
   - It then calls the `client.Register()` function to send a user registration request to the server.
   - If successful, the registration response is then marshaled to JSON, and the result is printed in green color.

4. **loginCmd:**
   - `loginCmd` is a subcommand that represents the `login` functionality of the CLI.
   - When invoked, it sets up a gRPC client (`grpcClient`) for communication with the server.
   - The `client.SetupGRPCClient()` function is used to set up the gRPC client.
   - It then calls the `client.LogIn()` function to send a user login request to the server.
   - If successful, the login response is then marshaled to JSON, and the result is printed in green color.


