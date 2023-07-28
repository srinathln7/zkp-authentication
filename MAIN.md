# Package `main`:

The `main` package serves as the main entry point for the ZKP Authentication CLI application. It initializes the command-line flags, starts the gRPC server in the background (if the `--server` flag is specified), and executes Cobra commands otherwise to handle CLI subcommands. The application gracefully shuts down the server upon receiving a shutdown signal and logs any errors that occur during execution. 

1. **init Function:**
   - The `init` function is used to set up the command-line flags for the application.
   - It calls the `SetupFlags` function from the `cmd` package, which sets up the flags for user and password.

2. **main Function:**
   - The `main` function is the main entry point of the ZKP Authentication CLI application.
   - It first defines a flag `runServerInBackground` using the `flag` package. This flag is used to determine whether to run the gRPC server in the background.
   - It then calls `flag.Parse()` to parse the command-line flags and set their values accordingly.

3. **Running gRPC Server:**
   - The application checks if the `--server` flag is set (`*runServerInBackground` is `true`).
   - If the `--server` flag is set, it indicates that the gRPC server should be started in the background.
   - It generates Chaum-Pedersen Zero-Knowledge Proof (CPZKP) system parameters by calling `cp_zkp.NewCPZKP()`.
   - If an error occurs during system parameter generation, it logs the error and exits the application with an error code.
   - If the system parameters are generated successfully, it creates a server configuration `cfg` with the CPZKP parameters.
   - It starts the gRPC server in the background by calling `server.RunServer(cfg)` inside a goroutine.

4. **Graceful Shutdown:**
   - After starting the gRPC server, the application waits for a graceful shutdown signal (e.g., Ctrl+C) using `signal.Notify`.
   - It creates a channel `c` to receive the signal.
   - The application blocks until it receives a signal on the channel `c`, indicating the user wants to shut down the server.
   - When the shutdown signal is received, the goroutine running the server stops, and the application exits.

5. **Cobra Command Execution:**
   - If the `--server` flag is not set (indicating that the server should not be run in the background), the application proceeds to execute the Cobra commands.
   - It calls `cmd.RootCmd.Execute()` to execute the root Cobra command, which is responsible for handling the CLI subcommands.

6. **Exit on Error:**
   - If an error occurs during command execution (e.g., invalid command or arguments), the application logs the error and exits with an error code.



