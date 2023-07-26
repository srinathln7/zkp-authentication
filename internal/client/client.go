package client

import (
	"context"
	"log"
	"math/big"
	"net"

	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	"github.com/srinathLN7/zkp_auth/lib/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	sys_config "github.com/srinathLN7/zkp_auth/lib/config"
)

type Client struct{}

type RegRes struct {
	Msg string `json:"msg"`
}

type LogInRes struct {
	SessionId string `json:"session_id"`
}

func RunClient(user, password string) {

	listener, err := net.Listen("tcp", ":"+sys_config.GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}

	// Set up the gRPC client
	grpcClientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(listener.Addr().String(), grpcClientOptions...)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create the gRPC client
	grpcClient := api.NewAuthClient(conn)

	// Run your client call logics

	client := &Client{}

	cpzkpParams := client.getCPZKPParams()
	x := client.getSecretValue(password)

	client.Register(grpcClient, user, cpzkpParams, x)
	client.LogIn(grpcClient, user, cpzkpParams, x)
}

// RegisterUser Registers the user with the given password and returns a message, if successful
func (c *Client) Register(grpcClient api.AuthClient, user string, cpzkpParams *cp_zkp.CPZKPParams, x *big.Int) (*RegRes, error) {
	ctx := context.Background()

	// Create a new Prover (Client) based on the generated secret value `x`
	client := cp_zkp.NewProver(x)

	// Prover(client) generates y1 and y2 values
	y1, y2 := client.GenerateYValues(cpzkpParams)

	// Received response
	_, err := grpcClient.Register(
		ctx,
		&api.RegisterRequest{
			User: user,
			Y1:   y1.String(),
			Y2:   y2.String(),
		},
	)

	if err != nil {
		log.Fatal(err)
		return nil, grpc_err.ErrRegistrationFailed
	}

	return &RegRes{
		Msg: "user registration successful",
	}, nil
}

// LogIn : Validates the login credentials using the Chaum-Pedersen Zero-Knowledge Proof
// protocol and returns a succesful message for a valid login
func (c *Client) LogIn(grpcClient api.AuthClient, user string, cpzkpParams *cp_zkp.CPZKPParams, x *big.Int) (*LogInRes, error) {
	ctx := context.Background()

	// Create a new Prover (Client) based on the generated secret value `x`
	client := cp_zkp.NewProver(x)

	k, r1, r2, err := client.CreateProofCommitment(cpzkpParams)
	if err != nil {
		log.Fatal(err)
		return nil, grpc_err.ErrLoginFailed
	}

	recvAuthChallengeRes, err := grpcClient.CreateAuthenticationChallenge(
		ctx,
		&api.AuthenticationChallengeRequest{
			User: user,
			R1:   r1.String(),
			R2:   r2.String(),
		},
	)

	if err != nil {
		log.Fatal(err)
		return nil, grpc_err.ErrLoginFailed
	}

	authID := recvAuthChallengeRes.AuthId
	cStr := recvAuthChallengeRes.C
	C, err := util.ParseBigInt(cStr, "c")
	if err != nil {
		log.Fatal(err)
		return nil, grpc_err.ErrLoginFailed
	}

	// Challenge response
	s := client.CreateProofChallengeResponse(k, C, cpzkpParams)

	// Verification Step
	verifyRes, err := grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s.String(),
		},
	)

	if err != nil {
		log.Fatal(err)
		return nil, grpc_err.ErrLoginFailed
	}

	return &LogInRes{
		SessionId: verifyRes.SessionId,
	}, nil
}

// getCPZKPParams generates the system parameters
func (c *Client) getCPZKPParams() *cp_zkp.CPZKPParams {
	// Generate the system parameters
	cpzkp, err := cp_zkp.NewCPZKP()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	cpzkpParams, err := cpzkp.InitCPZKPParams()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return cpzkpParams
}

// getSecretValue gets the secret value `x` by converting the password uniquely to big Int
func (c *Client) getSecretValue(password string) *big.Int {
	// Get the secret value `x` by converting the password uniquely to big Int
	return util.StringToUniqueBigInt(password)
}
