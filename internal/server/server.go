package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/google/uuid"
	grpc_err "github.com/srinathLN7/zkp_auth/api/v2/err"
	api "github.com/srinathLN7/zkp_auth/api/v2/proto"
	cp_zkp "github.com/srinathLN7/zkp_auth/internal/cpzkp"
	"google.golang.org/grpc"
)

type CPZKP interface {
	InitCPZKPParams() (*cp_zkp.CPZKPParams, error)
	SetZKPParams(p, q, g, h *big.Int) (*cp_zkp.CPZKPParams, error)
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

	if _, userExists := s.RegDir[req.User]; userExists {
		return nil, fmt.Errorf("user %s is already registered on the server", req.User)
	}

	var Y1, Y2 *big.Int = new(big.Int), new(big.Int)

	log.Printf("[grpc_server] req.Y1 %s", req.Y1)
	Y1, validY1 := Y1.SetString(req.Y1, 10)
	if !validY1 {
		return nil, errors.New("error parsing string Y1 to big.Int")
	}

	log.Printf("[grpc_server] req.Y2 %v", req.Y2)
	Y2, validY2 := Y2.SetString(req.Y2, 10)
	if !validY2 {
		return nil, errors.New("error parsing string Y2 to big.Int")
	}

	log.Printf("[grpc_server] Y1 %v", Y1)
	log.Printf("[grpc_server] Y2 %v", Y2)

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
	cpzkpParams, err = s.Config.CPZKP.SetZKPParams(big.NewInt(23), big.NewInt(11), big.NewInt(4), big.NewInt(9))
	if err != nil {
		return nil, err
	}

	// Create a verifier to create the authentication challenge
	verifier := &cp_zkp.Verifier{}
	c, err := verifier.CreateProofChallenge(cpzkpParams)
	if err != nil {
		return nil, err
	}

	c = big.NewInt(4)

	// Store the generated random value `c` and the `auth_id` in the authentication directory
	// for authentication verification process in the next step

	// We use the google's widely used `uuid` pkg to generate the authID
	authID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	var R1, R2 *big.Int = new(big.Int), new(big.Int)

	R1, validR1 := R1.SetString(req.R1, 10)
	if !validR1 {
		return nil, errors.New("error parsing string R1 to big.Int")
	}

	R2, validR2 := R2.SetString(req.R2, 10)
	if !validR2 {
		return nil, errors.New("error parsing string R2 to big.Int")
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
	cpzkpParams, err = s.Config.CPZKP.SetZKPParams(big.NewInt(23), big.NewInt(11), big.NewInt(4), big.NewInt(9))
	if err != nil {
		return nil, err
	}

	// Get the user name and `c` from the current `auth_id`

	log.Printf("server Auth Directory map %#v\n", s.AuthDir)
	user := s.AuthDir[req.AuthId].user
	c := s.AuthDir[req.AuthId].c

	// Retrieve y1 and y2

	log.Printf("server Registry directory map %#v\n", s.RegDir)

	y1 := s.RegDir[user].y1
	y2 := s.RegDir[user].y2

	// Create a prover to calculate the r1 and r2 values as part of the commitment step in the proof
	// prover := &cp_zkp.Prover{}
	// _, r1, r2, err := prover.CreateProofCommitment(cpzkpParams)
	// if err != nil {
	// 	return nil, err
	// }

	// r1 = big.NewInt(8)
	// r2 = big.NewInt(4)

	r1 := s.AuthDir[req.AuthId].r1
	r2 := s.AuthDir[req.AuthId].r2

	// convert `req.S` to big.Int
	var s_zkp *big.Int = new(big.Int)

	_, validS := s_zkp.SetString(req.S, 10)
	if !validS {
		return nil, errors.New("error parsing string `S` to big.Int")
	}

	// Create a verifier to verify the challenge
	verifier := &cp_zkp.Verifier{}
	isValidProof := verifier.VerifyProof(y1, y2, r1, r2, c, s_zkp, cpzkpParams)
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
