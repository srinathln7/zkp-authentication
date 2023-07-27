package server

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"github.com/srinathLN7/zkp_auth/lib/util"
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
	r1   *big.Int
	r2   *big.Int
}

type grpcServer struct {
	api.UnimplementedAuthServer

	// Simulate a server-side user directory
	// store the `y1` and `y2` of the specific user
	// Limited by in-memory non-persistence storage
	RegDir map[string]RegParams

	// Simulate a server-side authentication directory
	// store the `c`, `y1` and `y2` of the specific user
	// Limited by in-memory non-persistence storage
	AuthDir map[string]AuthParams

	*Config
}

func RunServer(config *Config) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
		return
	}

	grpcServerAddr := os.Getenv("SERVER_ADDRESS")
	listener, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
		return
	}

	// Create a new gRPC server and register the service
	grpcServer, err := NewGRPCSever(config)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	// Listen on the specified grpc server port

	log.Printf("grpc server listening on: %s\n", listener.Addr().String())

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}

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

	if _, userExists := s.RegDir[req.User]; userExists {
		return nil, grpc_err.ErrInvalidRegistration{User: req.User}
	}

	Y1, err := util.ParseBigInt(req.Y1, "y1")
	if err != nil {
		return nil, err
	}

	Y2, err := util.ParseBigInt(req.Y2, "y2")
	if err != nil {
		return nil, err
	}

	s.RegDir[req.User] = RegParams{
		y1: Y1,
		y2: Y2,
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

	// We use the google's widely used `uuid` pkg to generate the authID
	authID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	R1, err := util.ParseBigInt(req.R1, "r1")
	if err != nil {
		return nil, err
	}

	R2, err := util.ParseBigInt(req.R2, "r2")
	if err != nil {
		return nil, err
	}

	auth_id := authID.String()
	s.AuthDir[auth_id] = AuthParams{user: req.User,
		c:  c,
		r1: R1,
		r2: R2,
	}

	return &api.AuthenticationChallengeResponse{
		AuthId: auth_id,
		C:      c.String(),
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

	// Retrieve y1, y2, r1, r2
	y1 := s.RegDir[user].y1
	y2 := s.RegDir[user].y2
	r1 := s.AuthDir[req.AuthId].r1
	r2 := s.AuthDir[req.AuthId].r2

	// convert `req.S` to big.Int
	S, err := util.ParseBigInt(req.S, "s")
	if err != nil {
		return nil, err
	}

	// Create a verifier to verify the challenge
	verifier := &cp_zkp.Verifier{}
	isValidProof := verifier.VerifyProof(y1, y2, r1, r2, c, S, cpzkpParams)
	if !isValidProof {
		return nil, grpc_err.ErrInvalidChallengeResponse{S: req.S}
	}

	// If a valid proof is presented - then generate a sessionID and pass it as a response
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &api.AuthenticationAnswerResponse{SessionId: sessionID.String()}, nil
}
