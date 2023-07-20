package cp_zkp

import (
	"crypto/rand"
	"math/big"
)

// ZKPParameters represents the public parameters for the ZKP protocol.
type ZKPParameters struct {
	p, q, g, h *big.Int
}

// Prover represents the prover in the ZKP protocol.
type Prover struct {
	x *big.Int // Secret number x
}

// Verifier represents the verifier in the ZKP protocol.
type Verifier struct {
}

// NewZKPParameters generates random public parameters for the ZKP protocol.
// The function generates secure prime number p and finds a suitable prime number q based on p.
// It then calculates g and h based on p and q.
// The system parameters are generated as per the explanation in the forum: https://crypto.stackexchange.com/questions/99262/chaum-pedersen-protocol
func NewZKPParameters() (*ZKPParameters, error) {
	// Use larger prime number (e.g., 512 bits) for p
	p, err := rand.Prime(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	// Find a suitable prime q that divides (p-1) evenly
	q, err := findSuitableQ(p)
	if err != nil {
		return nil, err
	}

	// Calculate g = 2^(p-1)/q mod p
	g := new(big.Int).Exp(big.NewInt(2), new(big.Int).Div(new(big.Int).Sub(p, big.NewInt(1)), q), p)

	// Calculate h = random value in the range [1, p-1]
	h, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(1)))
	if err != nil {
		return nil, err
	}
	h.Add(h, big.NewInt(1))

	return &ZKPParameters{
		p: p,
		q: q,
		g: g,
		h: h,
	}, nil
}

// findSuitableQ finds a suitable prime q that divides (p-1) evenly.
func findSuitableQ(p *big.Int) (*big.Int, error) {
	for {
		q, err := rand.Prime(rand.Reader, 256) // Use larger prime number (e.g., 256 bits) for q
		if err != nil {
			return nil, err
		}

		if new(big.Int).Mod(new(big.Int).Sub(p, big.NewInt(1)), q).Cmp(big.NewInt(0)) == 0 {
			return q, nil
		}
	}
}

// NewProver creates a new Prover with the given secret password x.
func NewProver(x *big.Int) *Prover {
	return &Prover{
		x: x,
	}
}

// GenerateYValues generates y1 and y2 for the prover based on the public parameters.
// The prover calculates y1 = g^x mod p and y2 = h^x mod p.
func (p *Prover) GenerateYValues(params *ZKPParameters) (y1, y2 *big.Int) {
	y1 = new(big.Int).Exp(params.g, p.x, params.p)
	y2 = new(big.Int).Exp(params.h, p.x, params.p)
	return y1, y2
}

// CreateProof creates a zero-knowledge proof based on the prover's y1 and y2 values.
// The prover selects a random value k and commits (r1, r2) = (g^k mod p, h^k mod p).
// Then, the prover computes c = y1^k mod p.
// Finally, the prover calculates s = (k - c * x) mod q.
func (p *Prover) CreateProof(y1, y2 *big.Int, params *ZKPParameters) (r1, r2, c, s *big.Int, err error) {
	k, err := rand.Int(rand.Reader, params.q) // Use cryptographically secure random number generator
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Compute commitments (r1, r2) = (g^k mod p, h^k mod p)
	r1 = new(big.Int).Exp(params.g, k, params.p)
	r2 = new(big.Int).Exp(params.h, k, params.p)

	// Compute c = y1^k mod p
	c = new(big.Int).Exp(y1, k, params.p)

	// Compute s = (k - c * x) mod q
	s = new(big.Int).Sub(k, new(big.Int).Mul(c, p.x))
	s.Mod(s, params.q)

	return r1, r2, c, s, nil
}

// VerifyProof verifies the zero-knowledge proof using the verifier's y1, y2, and the public parameters.
// The verifier checks if r1 = (g^s * y1^c) mod p and r2 = (h^s * y2^c) mod p.
// If both checks pass, the proof is valid, and the function returns true; otherwise, it returns false.
func (v *Verifier) VerifyProof(y1, y2, r1, r2, c, s *big.Int, params *ZKPParameters) bool {
	left1 := new(big.Int).Exp(params.g, s, params.p)
	left1.Mul(left1, new(big.Int).Exp(y1, c, params.p))
	left1.Mod(left1, params.p)

	right1 := new(big.Int).Exp(r1, params.p, params.p)
	if left1.Cmp(right1) != 0 {
		return false
	}

	left2 := new(big.Int).Exp(params.h, s, params.p)
	left2.Mul(left2, new(big.Int).Exp(y2, c, params.p))
	left2.Mod(left2, params.p)

	right2 := new(big.Int).Exp(r2, params.p, params.p)
	return left2.Cmp(right2) == 0
}
