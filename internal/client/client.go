package client

import (
	"net"
	"testing"

	api "github.com/srinathLN7/zkp_auth/api/v1"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/internal/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGRPCServer(t *testing.T) {
	for sceanario, fn := range map[string]func(
		t *testing.T,
		grpcClient api.AuthClient,
		config *server.Config,
	){
		"register user succesfully":                    testRegisterUserSuccess,
		"create authentication challenge successfully": testAuthenticationChallengeSuccess,
		"verify authentication challenge successfully": testVerifyAuthenticationChallengeSuccess,
		"register existing user failure":               testRegisterUserFail,
		"create authentication challenge failure":      testAuthenticationChallengeFailure,
		"verify authentication challenge failure":      testVerifyAuthenticationChallengeFailure,
	} {
		t.Run(sceanario, func(t *testing.T) {
			grpcClient, config, teardown := setupGRPCClient(t, nil)
			defer teardown()
			fn(t, grpcClient, config)
		})
	}
}

func setupGRPCClient(t *testing.T, fn func(*server.Config)) (
	grpcClient api.AuthClient,
	cfg *server.Config,
	teardown func(),
) {
	// Helper marks the calling function as a test helper function. When printing file and line information, that function will be skipped
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	grpcClientOptions := []grpc.DialOption{grpc.WithInsecure()}
	cc, err := grpc.Dial(listener.Addr().String(), grpcClientOptions...)
	require.NoError(t, err)

	cpzkpParams, err := cp_zkp.NewCPZKP()
	require.NoError(t, err)

	cfg = &server.Config{
		CPZKP: cpzkpParams,
	}

	if fn != nil {
		fn(cfg)
	}

	grpcServer, err := server.NewGRPCSever(cfg)
	require.NoError(t, err)

	go func() {
		grpcServer.Serve(listener)
	}()

	grpcClient = api.NewAuthClient(cc)

	return grpcClient, cfg, func() {
		grpcServer.Stop()
		cc.Close()
		listener.Close()
	}
}
