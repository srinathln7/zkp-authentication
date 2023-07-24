package server

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/google/uuid"
	api "github.com/srinathLN7/zkp_auth/api/v1"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"google.golang.org/grpc"
)

type CPZKP interface {
	InitCPZKPParams() (*cp_zkp.CPZKPParams, error)
}

type Config struct {
	CPZKP CPZKP
}

type RegParams struct {
	y1 *big.Int
	y2 *big.Int
}

type AuthParams struct {
	user string
	c    *big.Int
}

type grpcServer struct {
	api.UnimplementedAuthServer

	// Simulate a server-side user directory
	// store the `y1` and `y2` of the specific user
	// Limited by in-memory non-persistence storage
	RegDir map[string]RegParams

	// Authentication id directory
	AuthDir map[string]AuthParams

	*Config
}

var _ api.AuthServer = (*grpcServer)(nil)

func newgrpcServer(config *Config) (*grpcServer, error) {
	// initialize the server with ZKP system params and an empty user directory
	return &grpcServer{
		RegDir:  make(map[string]RegParams),
		AuthDir: make(map[string]AuthParams),
		Config:  config,
	}, nil
}

// NewGRPCServer: creates a grpc server and registers the service to that server
func NewGRPCSever(config *Config) (*grpc.Server, error) {
	gsrv := grpc.NewServer()
	srv, err := newgrpcServer(config)
	if err != nil {
		return nil, err
	}
	api.RegisterAuthServer(gsrv, srv)
	return gsrv, nil
}

// Register: Simply registers a new grpc client (prover) on the server side
// by storing the passed-in req body containing `y1` and `y2` values
func (s *grpcServer) Register(ctx context.Context, req *api.RegisterRequest) (
	*api.RegisterResponse, error) {

	// ASSUMPTION: The `req.user` passed in for every user is UNIQUE
	// Check if the user already exists

	log.Println(s.RegDir[req.User])
	if _, userExists := s.RegDir[req.User]; userExists {
		return nil, fmt.Errorf("user %s is already registered on the server", req.User)
	}

	s.RegDir[req.User] = RegParams{
		y1: big.NewInt(req.Y1),
		y2: big.NewInt(req.Y2),
	}

	return &api.RegisterResponse{}, nil
}

func (s *grpcServer) CreateAuthenticationChallenge(ctx context.Context, req *api.AuthenticationChallengeRequest) (
	*api.AuthenticationChallengeResponse, error) {

	// First check if the user is registered on the server
	// Otherwise throw an error before proceeding further
	if _, userExists := s.RegDir[req.User]; !userExists {
		return nil, fmt.Errorf("user %s is not registered on the server", req.User)
	}

	// We use the google's widely used `uuid` pkg to generate the authID
	authID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	auth_id := authID.String()

	cpzkpParams, err := s.Config.CPZKP.InitCPZKPParams()
	if err != nil {
		return nil, err
	}

	// Create a verifier to create the authentication challenge
	verifier := &cp_zkp.Verifier{}
	c, err := verifier.CreateProofChallenge(cpzkpParams)
	if err != nil {
		return nil, err
	}

	// Store the generated random value `c` and the `auth_id` in the authentication directory
	// for authentication verification process in the next step
	s.AuthDir[auth_id] = AuthParams{user: req.User, c: c}

	return &api.AuthenticationChallengeResponse{
		AuthId: auth_id,
		C:      c.Int64(),
	}, nil
}

func (s *grpcServer) VerifyAuthentication(ctx context.Context, req *api.AuthenticationAnswerRequest) (
	*api.AuthenticationAnswerResponse, error) {

	// First check if the authentication id passed is valid
	if _, idExists := s.AuthDir[req.AuthId]; !idExists {
		return nil, fmt.Errorf("invalid authentication id: %s specified", req.AuthId)
	}

	// To verify the proof, we need the system params and
	// y1, y2, r1,r2, c, s
	cpzkpParams, err := s.Config.CPZKP.InitCPZKPParams()
	if err != nil {
		return nil, err
	}

	// Get the user name and `c` from the current `auth_id`
	user := s.AuthDir[req.AuthId].user
	c := s.AuthDir[req.AuthId].c

	// Retrieve y1 and y2
	y1 := s.RegDir[user].y1
	y2 := s.RegDir[user].y2

	// Create a prover to calculate the r1 and r2 values as part of the commitment step in the proof
	prover := &cp_zkp.Prover{}
	_, r1, r2, err := prover.CreateProofCommitment(cpzkpParams)
	if err != nil {
		return nil, err
	}

	// Create a verifier to verify the challenge
	verifier := &cp_zkp.Verifier{}
	isValidProof := verifier.VerifyProof(y1, y2, r1, r2, c, big.NewInt(req.S), cpzkpParams)
	if !isValidProof {
		return nil, api.ErrInvalidChallengeResponse{S: req.S}
	}

	// If a valid proof is presented - then generate a sessionID and pass it as a response
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &api.AuthenticationAnswerResponse{SessionId: sessionID.String()}, nil
}
