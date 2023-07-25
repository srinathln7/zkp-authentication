package server_test

import (
	"os"
	"testing"

	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	grpc_client "github.com/srinathLN7/zkp_auth/internal/client"
	grpc_server "github.com/srinathLN7/zkp_auth/internal/server"
)

// Run the tests
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGRPCServer(t *testing.T) {

	// We need a shared-server instance to mimic persistent storage during the test runtime
	// Hence, we setup the grpc client before running each of the individual test cases
	grpcClient, config, teardown := grpc_client.SetupGRPCClient(t, nil)

	// gracefully shutdown the server after running all the test cases
	defer teardown()

	for sceanario, fn := range map[string]func(
		t *testing.T,
		grpcClient api.AuthClient,
		config *grpc_server.Config,
	){
		"register user succesfully":       grpc_client.ClientRegisterUserSuccess,
		"register user failure":           grpc_client.ClientRegisterUserFail,
		"verification proof successfully": grpc_client.ClientVerifyProofSuccess,
		//"verification proof failure":      grpc_client.ClientVerifyProofFail,
	} {
		t.Run(sceanario, func(t *testing.T) {
			fn(t, grpcClient, config)
		})
	}
}
