package cp_zkp

import (
	"math/big"
	"testing"
)

// TestCPZKPProtocol quickly tests the correctness of the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol.
// Generates a prime number `p` of length 256 bits instead of 512 bits
// TestCPZKPProtocolFast quickly tests the correctness of the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol.
// Generates a prime number `p` of length 256 bits instead of 512 bits
func TestCPZKPProtocol(t *testing.T) {

	// Test multiple system parameters with different values of p and q

	// set bit size
	// n := 16
	// params1, err := InitZKPParameters(n) // Use a larger value for n (e.g., 256 bits)
	// if err != nil {
	// 	t.Errorf("Error generating ZKP parameters: %v", err)
	// 	return
	// }

	params1 := &ZKPParameters{
		p: big.NewInt(23),
		q: big.NewInt(11),
		g: big.NewInt(4),
		h: big.NewInt(9),
	}

	// Test if the protocol correctly invalidates a wrong proof and validates the right proof

	// Prover's secret value x
	x := big.NewInt(6)

	// Prover creates proof commitment (r1, r2)
	prover := NewProver(x)
	y1, y2 := prover.GenerateYValues(params1)
	k, r1, r2, err := prover.CreateProofCommitment(params1)
	if err != nil {
		t.Errorf("Error creating proof commitment: %v", err)
		return
	}

	// Verifier creates a challenge c
	verifier := Verifier{}
	c, err := verifier.CreateProofChallenge(params1)
	if err != nil {
		t.Errorf("Error creating challenge: %v", err)
		return
	}

	// Prover responds to the challenge and generates s
	s := prover.CreateProofChallengeResponse(k, c, params1)

	// Debug information
	t.Log("cp-zkp params: p=", params1.p, " q=", params1.q, " g=", params1.g, " h=", params1.h)
	t.Logf("y1: %v", y1)
	t.Logf("y2: %v", y2)
	t.Logf("r1: %v", r1)
	t.Logf("r2: %v", r2)
	t.Logf("c: %v", c)
	t.Logf("s: %v", s)

	// Calculate (g^s * y1^c) mod p
	exp1 := new(big.Int).Exp(y1, c, params1.p)
	tmp1 := new(big.Int).Mul(exp1, new(big.Int).Exp(params1.g, s, params1.p))
	tmp1.Mod(tmp1, params1.p)

	t.Logf("(g^s * y1^c) mod p: %v", tmp1)

	exp2 := new(big.Int).Exp(y2, c, params1.p)
	tmp2 := new(big.Int).Mul(new(big.Int).Exp(params1.h, s, params1.p), exp2)
	tmp2.Mod(tmp1, params1.p)

	t.Logf("(h^s * y2^c) mod p: %v", tmp2)

	// Test Correctness

	// Verifier verifies the proof
	valid := verifier.VerifyProof(y1, y2, r1, r2, c, s, params1)
	if !valid {
		t.Errorf("Proof validation failed: Expected valid proof, got invalid")
		return
	}

	// Test Soundness

	// Create an invalid proof by changing r1
	invalidR1 := new(big.Int).Add(r1, big.NewInt(1))

	// Verifier discards the invalid proof
	invalid := verifier.VerifyProof(y1, y2, invalidR1, r2, c, s, params1)
	if invalid {
		t.Errorf("Proof validation failed: Expected invalid proof, got valid")
		return
	}
}

// Run the tests
func TestMain(m *testing.M) {
	// Add any setup code here if needed
	// Run the tests
	m.Run()
}
