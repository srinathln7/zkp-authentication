package cp_zkp

import (
	"math/big"
	"testing"

	sys_config "github.com/srinathLN7/zkp_auth/lib/config"
	"github.com/srinathLN7/zkp_auth/lib/util"
)

// TestCPZKPProtocol tests the correctness and soundness of the
// Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol.
func TestCPZKPProtocol(t *testing.T) {

	cpZKP := &CPZKP{}
	params1, err := cpZKP.InitCPZKPParams()
	if err != nil {
		t.Errorf("error generating ZKP parameters: %v", err)
		return
	}

	// Test if the protocol correctly invalidates a wrong proof and validates the right proof

	// Prover's secret value x
	x, err := util.ParseBigInt(sys_config.CPZKP_TEST_X_CORRECT, "x")
	if err != nil {
		t.Errorf("error parsing the secret value `x` to big integer")
	}

	// Prover creates proof commitment (r1, r2)
	prover := NewProver(x)
	y1, y2 := prover.GenerateYValues(params1)
	k, r1, r2, err := prover.CreateProofCommitment(params1)
	if err != nil {
		t.Errorf("error creating proof commitment: %v", err)
		return
	}

	// Verifier creates a challenge c
	verifier := Verifier{}
	c, err := verifier.CreateProofChallenge(params1)
	if err != nil {
		t.Errorf("error creating challenge: %v", err)
		return
	}

	// Prover responds to the challenge and generates s
	s := prover.CreateProofChallengeResponse(k, c, params1)

	// Test Correctness

	// Verifier verifies the proof
	valid := verifier.VerifyProof(y1, y2, r1, r2, c, s, params1)
	if !valid {
		t.Errorf("proof validation failed: expected valid proof, got invalid")
		return
	}

	// Test Soundness

	// Create an invalid response to the verifier's  challenge
	invalidS := new(big.Int).Add(s, big.NewInt(1))

	// Verifier discards the invalid proof
	invalid := verifier.VerifyProof(y1, y2, r1, r2, c, invalidS, params1)
	if invalid {
		t.Errorf("proof validation failed: expected invalid proof, got valid")
		return
	}
}

// Run the tests
func TestMain(m *testing.M) {
	m.Run()
}
