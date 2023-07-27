package cmd

import (
	"encoding/json"
	"fmt"
	"log"

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
		fmt.Println("Use 'register' or 'login' subcommands.")
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		grpcClient, err := client.SetupGRPCClient()
		if err != nil {
			log.Fatalf("error setting up grpc client %s", err.Error())
			return
		}
		regRes, err := client.Register(*grpcClient, user, password)
		if err != nil {
			color.Red(err.Error())
			return
		}

		resJSON, err := json.Marshal(regRes)
		if err != nil {
			log.Fatal("error:", err)
			return
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
			return
		}
		loginRes, err := client.LogIn(*grpcClient, user, password)
		if err != nil {
			color.Red(err.Error())
			return
		}

		resJSON, err := json.Marshal(loginRes)
		if err != nil {
			log.Fatal("error:", err)
			return
		}

		color.Green(string(resJSON))
	},
}
