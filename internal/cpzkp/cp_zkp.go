package cp_zkp

import (
	"crypto/rand"
	"log"
	"math/big"

	"github.com/srinathLN7/zkp_auth/lib/config"
	"github.com/srinathLN7/zkp_auth/lib/util"
)

type CPZKP struct {
}

// CPZKPParams represents the public parameters for the ZKP protocol.
// p -> primeP, q -> primeQ, g -> generatorG, and h -> generatorH.
type CPZKPParams struct {
	p, q, g, h *big.Int
}

// Prover represents the prover in the ZKP protocol.
type Prover struct {
	x *big.Int // Secret number x
}

// Verifier represents the verifier in the ZKP protocol.
type Verifier struct {
}

func NewCPZKP() (*CPZKP, error) {
	return &CPZKP{}, nil
}

// InitCPZKPParams initializes the Chaum-Pedersen ZKP protocol system params.
func (zkp *CPZKP) InitCPZKPParams() (*CPZKPParams, error) {

	// `p` and `q` has 164 bits
	p, err := util.ParseBigInt(config.CPZKP_PARAM_P, "p")
	if err != nil {
		return nil, err
	}

	q, err := util.ParseBigInt(config.CPZKP_PARAM_Q, "q")
	if err != nil {
		return nil, err
	}

	g, err := util.ParseBigInt(config.CPZKP_PARAM_G, "g")
	if err != nil {
		return nil, err
	}

	h, err := util.ParseBigInt(config.CPZKP_PARAM_H, "h")
	if err != nil {
		return nil, err
	}

	return &CPZKPParams{
		p: p,
		q: q,
		g: g,
		h: h,
	}, nil
}

// NewProver creates a new Prover with the given secret password x.
func NewProver(x *big.Int) *Prover {
	return &Prover{
		x: x,
	}
}

// GenerateYValues generates y1 and y2 for the prover based on the public parameters.
// The prover calculates y1 = g^x mod p and y2 = h^x mod p.
// y1 and y2 are public informations
func (p *Prover) GenerateYValues(params *CPZKPParams) (y1, y2 *big.Int) {
	y1 = new(big.Int).Exp(params.g, p.x, params.p)
	y2 = new(big.Int).Exp(params.h, p.x, params.p)
	return y1, y2
}

// CreateProofCommitment: creates a zero-knowledge proof commitment step based on the prover's y1 and y2 values.
// The prover selects a random value k and commits (r1, r2) = (g^k mod p, h^k mod p).
func (p *Prover) CreateProofCommitment(params *CPZKPParams) (k, r1, r2 *big.Int, err error) {
	k, err = rand.Int(rand.Reader, params.q) // Use cryptographically secure random number generator
	if err != nil {
		return nil, nil, nil, err
	}

	// Compute commitments (r1, r2) = (g^k mod p, h^k mod p)
	r1 = new(big.Int).Exp(params.g, k, params.p)
	r2 = new(big.Int).Exp(params.h, k, params.p)

	return k, r1, r2, nil
}

// CreateProofChallenge: verifier creates a challenge to the prover by generating a random big integer
// `c` which will be subsequently used by the prover in the `CreateProofChallengeResponse` step
func (v *Verifier) CreateProofChallenge(params *CPZKPParams) (c *big.Int, err error) {
	// Generate a random `c` using cryptographically secure random number generator.
	c, err = rand.Int(rand.Reader, params.q)
	if err != nil {
		return nil, err
	}

	// Ensure c is not zero
	if c.Cmp(big.NewInt(0)) <= 0 {
		c.Add(c, big.NewInt(1))
	}

	return c, nil
}

// CreateProofChallengeResponse: prover creates the response to the verifier's challenge
// Compute s = (k - c * x) mod q
func (p *Prover) CreateProofChallengeResponse(k, c *big.Int, params *CPZKPParams) (s *big.Int) {
	s = new(big.Int).Sub(k, new(big.Int).Mul(c, p.x))
	s.Mod(s, params.q)

	if s.Cmp(big.NewInt(0)) < 0 {
		s.Add(s, params.q)
	}

	return s
}

// VerifyProof verifies the zero-knowledge proof using the verifier's y1, y2, and the public parameters.
// The verifier checks if r1 = (g^s * y1^c) mod p and r2 = (h^s * y2^c) mod p.
// If both checks pass, the proof is valid, and the function returns true; otherwise, it returns false.
func (v *Verifier) VerifyProof(y1, y2, r1, r2, c, s *big.Int, params *CPZKPParams) bool {

	// Debug logs
	log.Printf("[cp_zkp] y1 = %v", y1)
	log.Printf("[cp_zkp] y2 = %v", y2)
	log.Printf("[cp_zkp] r1 = %v", r1)
	log.Printf("[cp_zkp] r2 = %v", r2)
	log.Printf("[cp_zkp] c = %v", c)
	log.Printf("[cp_zkp] s = %v", s)

	// Remember: (ab) mod p = ( (a mod p) (b mod p)) mod p
	l1 := new(big.Int).Exp(params.g, s, params.p) // g^s .mod p
	l1.Mul(l1, new(big.Int).Exp(y1, c, params.p)) // y1^c mod p
	l1.Mod(l1, params.p)                          // (g^s mod p) (y1 ^c mod p) mod p = (g^s . y1^c). mod p

	if l1.Cmp(r1) != 0 {
		return false
	}

	l2 := new(big.Int).Exp(params.h, s, params.p)
	l2.Mul(l2, new(big.Int).Exp(y2, c, params.p))
	l2.Mod(l2, params.p)

	return l2.Cmp(r2) == 0
}
