package client

import (
	"context"
	"math/big"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/internal/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
)

// SetupGRPCClient: sets up the grpc client given the server config
func SetupGRPCClient(t *testing.T, fn func(*server.Config)) (
	grpcClient api.AuthClient,
	cfg *server.Config,
	teardown func(),
) {
	// Helper marks the calling function as a test helper function.
	// When printing file and line information, that function will be skipped
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	grpcClientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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

	grpServer, err := server.NewGRPCSever(cfg)
	require.NoError(t, err)

	go func() {
		grpServer.Serve(listener)
	}()

	grpcClient = api.NewAuthClient(cc)

	return grpcClient, cfg, func() {
		grpServer.Stop()
		cc.Close()
		listener.Close()
	}
}

func ClientRegisterUserSuccess(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Prover(client) generates y1 and y2 values
	// We set the secret value to `x=6`
	prover := cp_zkp.NewProver(big.NewInt(6))
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)
	y1, y2 := prover.GenerateYValues(cpzkpParams)

	// Desired response for successful user registration
	expectedResp := &api.RegisterResponse{}

	// Received response
	recvResp, err := grpcClient.Register(
		ctx,
		&api.RegisterRequest{
			User: "srinath",
			Y1:   y1.String(),
			Y2:   y2.String(),
		},
	)

	require.NoError(t, err)

	// Compare public fields of responses using cmp.Equal with IgnoreUnexported option
	// This is because a successful registration returns an empty response in our case
	// with only internally unexported generated fields
	if !cmp.Equal(recvResp, expectedResp, cmpopts.IgnoreUnexported(api.RegisterResponse{})) {
		t.Errorf("received response: %v, expected response: %v", recvResp, expectedResp)
	}
}

func ClientRegisterUserFail(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
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
			Y1:   y1.String(),
			Y2:   y2.String(),
		},
	)

	// We expect a non-nil error
	if err == nil {
		t.Fatal("expected a non-nil error")
	}

}

func ClientVerifyProofSuccess(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Verfication happens in two stages:

	// Step 1) Create Authentication Challenge

	t.Log("step 1: creating authentication challenge")

	// We set the correct secret value to `x=6`
	prover := cp_zkp.NewProver(big.NewInt(6))
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// Generate r1 and r2 values

	t.Log("prover creating proof commitment")
	k, r1, r2, err := prover.CreateProofCommitment(cpzkpParams)
	require.NoError(t, err)

	t.Log("cp-zkp params:", cpzkpParams)
	t.Logf("k: %v", k)
	t.Logf("r1: %v", r1)
	t.Logf("r2: %v", r2)

	// Create authentication challenge for user `srinath`

	t.Log("authentication challenge - connecting to grpc client")
	recvAuthChallengeRes, err := grpcClient.CreateAuthenticationChallenge(
		ctx,
		&api.AuthenticationChallengeRequest{
			User: "srinath",
			R1:   r1.String(),
			R2:   r2.String(),
		},
	)
	require.NoError(t, err)

	authID := recvAuthChallengeRes.AuthId
	cStr := recvAuthChallengeRes.C

	var c *big.Int = new(big.Int)
	_, validC := c.SetString(cStr, 10)
	require.True(t, validC)

	t.Logf("c: %v", c)

	// Step 2) Verify Authentication

	// Prover responds to the verifiers challenge

	t.Log("step 2: verify authentication")
	s := prover.CreateProofChallengeResponse(k, c, cpzkpParams)

	t.Logf("s: %v", s)

	// Create verification step
	t.Log("verify authentication - connecting to grpc client")
	_, err = grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s.String(),
		},
	)

	if err != nil {
		t.Fatal("verification failed for valid proof")
	}
}

func ClientVerifyProofFail(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
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
			R1:   r1.String(),
			R2:   r2.String(),
		},
	)
	require.NoError(t, err)

	authID := recvAuthChallengeRes.AuthId

	// Prover responds to the verifiers challenge incorrectly
	// Compute `s = (k - c * x) mod q`. Since prover has no knowledge of `x`, he cannot compute s correctly
	// Prover responds incorrectly to the verifiers challenge
	s := "555"

	// Create verification step
	_, err = grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s,
		},
	)

	// Expected err
	expErr := grpc_err.ErrInvalidChallengeResponse{S: s}
	if err != expErr {
		t.Fatalf("received err: %v, expected: %v", err, expErr)
	}

}
