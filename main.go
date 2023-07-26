package main

import (
	"log"
	"os"

	"github.com/srinathLN7/zkp_auth/cmd"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/internal/server"
)

func init() {
	cmd.SetupFlags()
}

func main() {
	cpzkpParams, err := cp_zkp.NewCPZKP()
	if err != nil {
		log.Println("error generating system parameters:", err)
		os.Exit(1)
	}

	cfg := &server.Config{
		CPZKP: cpzkpParams,
	}

	go server.RunServer(cfg)

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal("error:", err)
		os.Exit(1)
	}
}
