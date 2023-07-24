package server_test

import (
	"os"
	"testing"

	api "github.com/srinathLN7/zkp_auth/api/v1"
	grpc_client "github.com/srinathLN7/zkp_auth/internal/client"
	grpc_server "github.com/srinathLN7/zkp_auth/internal/server"
)

// Run the tests
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGRPCServer(t *testing.T) {
	for sceanario, fn := range map[string]func(
		t *testing.T,
		grpcClient api.AuthClient,
		config *grpc_server.Config,
	){
		"register user succesfully": grpc_client.ClientRegisterUserSuccess,
		//"register user failure":     grpc_client.ClientRegisterUserFail,
		//"verification proof successfully": grpc_client.ClientVerifyProofSuccess,
		//"verification proof failure":      grpc_client.ClientVerifyProofFail,
	} {
		t.Run(sceanario, func(t *testing.T) {
			grpcClient, config, teardown := grpc_client.SetupGRPCClient(t, nil)
			defer teardown()
			fn(t, grpcClient, config)
		})
	}
}
