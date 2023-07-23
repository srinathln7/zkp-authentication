package server

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"os"
	"testing"

	api "github.com/srinathLN7/zkp_auth/api/v1"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run the tests
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGRPCServer(t *testing.T) {
	for sceanario, fn := range map[string]func(
		t *testing.T,
		grpcClient api.AuthClient,
		config *Config,
	){
		"register user succesfully":       ClientRegisterUserSuccess,
		"register user failure":           ClientRegisterUserFail,
		"verification proof successfully": ClientVerifyProofSuccess,
		"verification proof failure":      ClientVerifyProofFail,
	} {
		t.Run(sceanario, func(t *testing.T) {
			grpcClient, config, teardown := setupGRPCClient(t, nil)
			defer teardown()
			fn(t, grpcClient, config)
		})
	}
}

func setupGRPCClient(t *testing.T, fn func(*Config)) (
	grpcClient api.AuthClient,
	cfg *Config,
	teardown func(),
) {
	// Helper marks the calling function as a test helper function. When printing file and line information, that function will be skipped
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	grpcClientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.Dial(listener.Addr().String(), grpcClientOptions...)
	require.NoError(t, err)

	cpzkpParams, err := cp_zkp.NewCPZKP()
	require.NoError(t, err)

	cfg = &Config{
		CPZKP: cpzkpParams,
	}

	if fn != nil {
		fn(cfg)
	}

	grpcServer, err := NewGRPCSever(cfg)
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

func ClientRegisterUserSuccess(t *testing.T, grpcClient api.AuthClient, config *Config) {
	ctx := context.Background()

	// Prover(client) generates y1 and y2 values
	// We set the secret value to `x=6`
	prover := cp_zkp.NewProver(big.NewInt(6))
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)
	y1, y2 := prover.GenerateYValues(cpzkpParams)

	// Desired response for successful user registration
	res := &api.RegisterResponse{}

	// Received response
	recvRes, err := grpcClient.Register(
		ctx,
		&api.RegisterRequest{
			User: "srinath",
			Y1:   y1.Int64(),
			Y2:   y2.Int64(),
		},
	)

	require.NoError(t, err)
	require.Equal(t, res, recvRes)
}

func ClientRegisterUserFail(t *testing.T, grpcClient api.AuthClient, config *Config) {
	ctx := context.Background()

	// Prover(client) generates y1 and y2 values
	prover := cp_zkp.NewProver(big.NewInt(6))
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// Generate y1 and y2 values
	y1, y2 := prover.GenerateYValues(cpzkpParams)

	// Register the user `srinath` again
	_, err = grpcClient.Register(
		ctx,
		&api.RegisterRequest{
			User: "srinath",
			Y1:   y1.Int64(),
			Y2:   y2.Int64(),
		},
	)

	// Desired error for failure user registration
	expErr := fmt.Errorf("user %s is already registered on the server", "srinath")
	if err != expErr {
		t.Fatalf("received err: %v, expected: %v", err, expErr)
	}
}

func ClientVerifyProofSuccess(t *testing.T, grpcClient api.AuthClient, config *Config) {
	ctx := context.Background()

	// Verfication happens in two stages:

	// Step 1) Create Authentication Challenge

	// We set the correct secret value to `x=6`
	prover := cp_zkp.NewProver(big.NewInt(6))
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// Generate r1 and r2 values
	k, r1, r2, err := prover.CreateProofCommitment(cpzkpParams)
	require.NoError(t, err)

	// Create authentication challenge for user `srinath`
	recvAuthChallengeRes, err := grpcClient.CreateAuthenticationChallenge(
		ctx,
		&api.AuthenticationChallengeRequest{
			User: "srinath",
			R1:   r1.Int64(),
			R2:   r2.Int64(),
		},
	)
	require.NoError(t, err)

	authID := recvAuthChallengeRes.AuthId
	c := recvAuthChallengeRes.C

	// Step 2) Verify Authentication

	// Prover responds to the verifiers challenge
	s := prover.CreateProofChallengeResponse(k, big.NewInt(c), cpzkpParams)

	// Create verification step
	_, err = grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s.Int64(),
		},
	)

	if err != nil {
		t.Fatal("verification failed for valid proof")
	}
}

func ClientVerifyProofFail(t *testing.T, grpcClient api.AuthClient, config *Config) {
	ctx := context.Background()

	// Verfication happens in two steps:
	// 1) Create Authentication Challenge
	// 2) Verify Authentication

	// Generate r1 and r2 values
	proverF := cp_zkp.Prover{}
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	_, r1, r2, err := proverF.CreateProofCommitment(cpzkpParams)
	require.NoError(t, err)

	// Create authentication challenge for user `srinath`
	recvAuthChallengeRes, err := grpcClient.CreateAuthenticationChallenge(
		ctx,
		&api.AuthenticationChallengeRequest{
			User: "srinath",
			R1:   r1.Int64(),
			R2:   r2.Int64(),
		},
	)
	require.NoError(t, err)

	authID := recvAuthChallengeRes.AuthId

	// Prover responds to the verifiers challenge incorrectly
	// Compute `s = (k - c * x) mod q`. Since prover has no knowledge of `x`, he cannot compute s correctly
	// Prover responds incorrectly to the verifiers challenge
	s := int64(555)

	// Create verification step
	_, err = grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s,
		},
	)

	// Expected err
	expErr := api.ErrInvalidChallengeResponse{S: s}
	if err != expErr {
		t.Fatalf("received err: %v, expected: %v", err, expErr)
	}

}
