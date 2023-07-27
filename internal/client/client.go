package client

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	"github.com/srinathLN7/zkp_auth/lib/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
)

type RegRes struct {
	Msg string `json:"msg"`
}

type LogInRes struct {
	SessionId string `json:"session_id"`
}

func SetupGRPCClient() (*api.AuthClient, error) {

	// Set up the gRPC client
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
		return nil, err
	}

	grpcServerAddr := os.Getenv("SERVER_ADDRESS")
	log.Printf("grpc client dialing on server address %s", grpcServerAddr)

	grpcClientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(grpcServerAddr, grpcClientOptions...)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
		return nil, err
	}

	// Create the gRPC client
	grpcClient := api.NewAuthClient(conn)
	return &grpcClient, nil
}

// RegisterUser Registers the user with the given password and returns a message, if successful
func Register(grpcClient api.AuthClient, user, password string) (*RegRes, error) {

	// Generate the system parameters
	cpzkp, err := cp_zkp.NewCPZKP()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cpzkpParams, err := cpzkp.InitCPZKPParams()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Get the secret value `x` by converting the password uniquely to big Int
	x := getSecretValue(password)
	log.Println("[grpcClient-Prover] Transformed password in to a secret value `x`")

	// Create a new Prover (Client) based on the generated secret value `x`
	// to calculate the y1 and y2 params
	client := cp_zkp.NewProver(x)

	// Prover(client) generates y1 and y2 values
	y1, y2 := client.GenerateYValues(cpzkpParams)

	// Received response
	ctx := context.Background()
	_, err = grpcClient.Register(
		ctx,
		&api.RegisterRequest{
			User: user,
			Y1:   y1.String(),
			Y2:   y2.String(),
		},
	)

	if err != nil {
		log.Fatal(color.RedString(err.Error()))
		return nil, grpc_err.ErrInvalidRegistration{User: user}
	}

	return &RegRes{
		Msg: " user registration successful ",
	}, nil
}

// LogIn : Validates the login credentials using the Chaum-Pedersen Zero-Knowledge Proof
// protocol and returns a succesful message for a valid login
func LogIn(grpcClient api.AuthClient, user, password string) (*LogInRes, error) {

	// Generate the system parameters
	cpzkp, err := cp_zkp.NewCPZKP()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cpzkpParams, err := cpzkp.InitCPZKPParams()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Get the secret value `x` by converting the password uniquely to big Int
	x := getSecretValue(password)

	log.Println("[grpcClient-Prover] Retrieved secret value `x` from the input password")

	// Create a new Prover (Client) based on the generated secret value `x`
	// to calculate the r1 and r2 params for committing the proof
	client := cp_zkp.NewProver(x)

	k, r1, r2, err := client.CreateProofCommitment(cpzkpParams)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctx := context.Background()
	recvAuthChallengeRes, err := grpcClient.CreateAuthenticationChallenge(
		ctx,
		&api.AuthenticationChallengeRequest{
			User: user,
			R1:   r1.String(),
			R2:   r2.String(),
		},
	)

	if err != nil {
		log.Fatal(color.RedString(err.Error()))
		return nil, err
	}

	authID := recvAuthChallengeRes.AuthId
	cStr := recvAuthChallengeRes.C
	c, err := util.ParseBigInt(cStr, "c")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Challenge response

	s := client.CreateProofChallengeResponse(k, c, cpzkpParams)

	// Verification Step
	verifyRes, err := grpcClient.VerifyAuthentication(
		ctx,
		&api.AuthenticationAnswerRequest{
			AuthId: authID,
			S:      s.String(),
		},
	)

	if err != nil {
		log.Fatal(color.RedString(err.Error()))
		return nil, grpc_err.ErrInvalidChallengeResponse{S: s.String()}
	}

	return &LogInRes{
		SessionId: verifyRes.SessionId,
	}, nil

}

// getSecretValue gets the secret value `x` by converting the password uniquely to big Int
func getSecretValue(password string) *big.Int {
	// Get the secret value `x` by converting the password uniquely to big Int
	return util.StringToUniqueBigInt(password)
}
