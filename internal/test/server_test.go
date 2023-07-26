package server_test

import (
	"os"
	"testing"

	grpc_client "github.com/srinathLN7/zkp_auth/internal/client"
)

// Run the tests
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGRPCServer(t *testing.T) {

	// We need a shared-server instance to mimic persistent storage during the test runtime
	// Hence, we setup the grpc client before running each of the individual test cases
	grpcClient, config, teardown := grpc_client.SetupGRPCClient(t, nil)

	// gracefully shutdown the server after finishing all the test cases
	defer teardown()

	// To ensure that the test cases are run sequentially, we use subtests.
	// Subtests are a way to group related tests together and run them in a specific order.
	// Each test case is run in its own subtest to provide better isolation, readability, and reporting.

	t.Run("register user succesfully", func(t *testing.T) {
		grpc_client.ClientRegisterUserSuccess(t, grpcClient, config)
	})

	t.Run("verification proof successfully", func(t *testing.T) {
		grpc_client.ClientVerifyProofSuccess(t, grpcClient, config)
	})

	t.Run("verification proof failure", func(t *testing.T) {
		grpc_client.ClientVerifyProofFail(t, grpcClient, config)
	})

	t.Run("register user failure", func(t *testing.T) {
		grpc_client.ClientRegisterUserFail(t, grpcClient, config)
	})

}
