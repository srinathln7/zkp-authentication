package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/srinathLN7/zkp_auth/cmd"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/internal/server"
)

func init() {
	cmd.SetupFlags()
}

func main() {

	var runServerInBackground = flag.Bool("server", false, "run grpc server in the background")
	flag.Parse()

	// Check if the --server flag is set
	if *runServerInBackground {
		cpzkpParams, err := cp_zkp.NewCPZKP()
		if err != nil {
			log.Fatal("error generating system parameters:", err)
			os.Exit(1)
		}

		cfg := &server.Config{
			CPZKP: cpzkpParams,
		}

		// Create and start the gRPC server in the background
		// To do this, we spin up a new go routine
		go server.RunServer(cfg)

		// Wait for a graceful shutdown signal (e.g., Ctrl+C)
		// This will keep the main function running and prevent it from exiting immediately
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		// If the server is running, return to prevent executing Cobra commands
		return
	}

	//  Execute the Cobra commands otherwise
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal("error:", err)
		os.Exit(1)
	}
}
