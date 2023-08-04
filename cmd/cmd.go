package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/srinathLN7/zkp_auth/internal/client"
)

var (
	user     string
	password string
)

func SetupFlags() {
	RootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "User")
	RootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password")
	RootCmd.AddCommand(registerCmd)
	RootCmd.AddCommand(loginCmd)
}

var RootCmd = &cobra.Command{
	Use:   "zkp_auth",
	Short: "A CLI for ZKP Authentication",
	Run: func(cmd *cobra.Command, args []string) {

		color.Yellow("************************ Welcome to Chaum-Pedersen ZKP based Authentication CLI *****************")
		color.Yellow("Please use 'register' or 'login' subcommands.")
		color.Yellow("To register a new user run: `go run main.go -register -u <username> -p <password>`")
		color.Yellow("To login a registered user run: `go run main.go -login -u <username> -p <password>`")
		color.Yellow("To exit this terminal press CTRL+C")

		// Setup a signal handler to capture interrupt and termination signals
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGTERM)

		// when done exit the program gracefully
		<-done

	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		grpcClient, err := client.SetupGRPCClient()
		if err != nil {
			log.Fatalf("error setting up grpc client %s", err.Error())
		}
		regRes, err := client.Register(*grpcClient, user, password)
		if err != nil {
			return
		}

		resJSON, err := json.Marshal(regRes)
		if err != nil {
			log.Fatal("error:", err)
		}

		color.Green(string(resJSON))
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with a registered user",
	Run: func(cmd *cobra.Command, args []string) {
		grpcClient, err := client.SetupGRPCClient()
		if err != nil {
			log.Fatalf("error setting up grpc client %s", err.Error())
		}
		loginRes, err := client.LogIn(*grpcClient, user, password)
		if err != nil {
			return
		}

		resJSON, err := json.Marshal(loginRes)
		if err != nil {
			log.Fatal("error:", err)
		}

		color.Green(string(resJSON))
	},
}
