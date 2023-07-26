package test

import (
	"context"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/internal/server"
	sys_config "github.com/srinathLN7/zkp_auth/lib/config"
	"github.com/srinathLN7/zkp_auth/lib/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	listener, err := net.Listen("tcp", ":"+sys_config.GRPC_PORT)
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

// ClientRegisterUserSuccess : Tests registering the client on the server successfull sceanario
func testClientRegisterUserSuccess(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Generate the system parameters
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// We use the CORRECT secret value here and register here
	x, err := util.ParseBigInt(sys_config.CPZKP_TEST_X_CORRECT, "x")
	require.NoError(t, err)
	prover := cp_zkp.NewProver(x)

	// Prover(client) generates y1 and y2 values
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

// ClientRegisterUserFail : Tests registering the client on the server failure sceanario
func testClientRegisterUserFail(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Generate the system parameters
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	x, err := util.ParseBigInt(sys_config.CPZKP_TEST_X_CORRECT, "x")
	require.NoError(t, err)
	prover := cp_zkp.NewProver(x)

	// Prover(client) generates y1 and y2 values
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

// ClientVerifyProofSuccess : Tests a client generating a valid proof sceanario
func testClientVerifyProofSuccess(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Generate the system parameters
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// Verfication happens in two stages:

	// Step 1) Create Authentication Challenge

	// We use the CORRECT secret value here for generating a successful proof
	// and test the success sceanario
	x, err := util.ParseBigInt(sys_config.CPZKP_TEST_X_CORRECT, "x")
	require.NoError(t, err)
	proverS := cp_zkp.NewProver(x)

	t.Log("step 1: creating authentication challenge")

	// Generate r1 and r2 values
	t.Log("prover creating proof commitment")
	k, r1, r2, err := proverS.CreateProofCommitment(cpzkpParams)
	require.NoError(t, err)

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
	c, err := util.ParseBigInt(cStr, "c")
	require.NoError(t, err)

	t.Logf("c: %v", c)

	// Step 2) Verify Authentication

	// Prover responds to the verifiers challenge

	t.Log("step 2: verify authentication")
	s := proverS.CreateProofChallengeResponse(k, c, cpzkpParams)

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
		t.Fatal("Verification failed for valid proof")
	}
}

// ClientVerifyProofFail : Tests a client generating a invalid proof sceanario
func testClientVerifyProofFail(t *testing.T, grpcClient api.AuthClient, config *server.Config) {
	ctx := context.Background()

	// Generate the system parameters
	cpzkpParams, err := config.CPZKP.InitCPZKPParams()
	require.NoError(t, err)

	// We use the INCORRECT secret value here for generating a failure proof
	// and test the failure sceanario
	x, err := util.ParseBigInt(sys_config.CPZKP_TEST_X_INCORRECT, "x")
	require.NoError(t, err)
	proverF := cp_zkp.NewProver(x)

	// Generate r1 and r2 values
	k, r1, r2, err := proverF.CreateProofCommitment(cpzkpParams)
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

	cStr := recvAuthChallengeRes.C
	c, err := util.ParseBigInt(cStr, "c")
	require.NoError(t, err)

	// Prover responds to the verifiers challenge incorrectly
	// Compute `s = (k - c * x) mod q`. Since prover has no knowledge of `x`, he cannot compute s correctly
	// Prover responds incorrectly to the verifiers challenge

	s := proverF.CreateProofChallengeResponse(k, c, cpzkpParams)

	// Create verification step
	_, err = grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s.String(),
		},
	)

	// Expected err
	expErr := grpc_err.ErrInvalidChallengeResponse{S: s.String()}

	// Check if both of them are equal
	require.Equal(t, expErr.Error(), err.Error())
}
