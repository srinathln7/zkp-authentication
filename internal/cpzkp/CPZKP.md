
# Package `cp_zkp` 

The CP-ZKP protocol is implemented using the `cp_zkp` package which consits of two files `cp_zkp.go` and `cp_zkp_test.go`


## Implementation

The `cp_zkp.go` file includes the necessary imports and data structures for the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol.

### Data Structures:

- `CPZKPParams` Struct: Holds the public parameters of the ZKP protocol, namely prime numbers `p` and `q`, and generator values `g` and `h`.

- `Prover` Struct: Represents the prover in the ZKP protocol containing the secret value `x`.

- `Verifier` Struct: Represents the verifier in the ZKP protocol.

### Functions and Methods:

- `NewCPZKP() (*CPZKP, error)`: Initializes and returns a new CPZKP instance.

- `InitCPZKPParams() (*CPZKPParams, error)`: Generates the system parameters `p`, `q`, `g`, and `h` from the configuration file. It logs the generated parameters to the console and returns them as a `CPZKPParams` struct.

- `NewProver(x *big.Int) *Prover`: Creates a new prover instance with the given secret value `x`.

- `GenerateYValues(params *CPZKPParams) (y1, y2 *big.Int)`: Calculates `y1 = g^x mod p` and `y2 = h^x mod p` based on the prover's secret value `x` and the public parameters. It logs the generated `y1` and `y2` values to the console and returns them.

- `CreateProofCommitment(params *CPZKPParams) (k, r1, r2 *big.Int, err error)`: Creates a proof commitment step. It selects a random value `k` from the range [0, `p`) and computes commitments `r1 = g^k mod p` and `r2 = h^k mod p`. It logs the generated `k`, `r1`, and `r2` values to the console and returns them. Note if `k` is found to be zero, then another random number is generated until `k` is not zero. This is a design choice to ensure that we are indeed working with big integers.

- `CreateProofChallenge(params *CPZKPParams) (c *big.Int, err error)`: Generates a random challenge `c` from the range [o, `p`) and logs the generated `c` value to the console and returns it. Note if `c` is found to be zero, then another random number is generated until `c` is not zero. This is a design choice to ensure that we are indeed working with big integers.

- `CreateProofChallengeResponse(k, c *big.Int, params *CPZKPParams) (s *big.Int)`: Calculates the prover's response `s` to the verifier's challenge `c`. It computes `s = (k - c * x) mod q` and logs the computed `s` value to the console and returns it.

- `VerifyProof(y1, y2, r1, r2, c, s *big.Int, params *CPZKPParams) bool`: Verifies the zero-knowledge proof using the verifier's values and the public parameters. It checks whether `r1 = (g^s * y1^c) mod p` and `r2 = (h^s * y2^c) mod p`. If both checks pass, the proof is valid, and the function returns `true`; otherwise, it returns `false`.

Overall, the CP-ZKP protocol allows a prover to demonstrate knowledge of a secret value `x` without revealing it to a verifier. The prover generates proof commitments `(r1, r2)` and responds to the verifier's challenge `s` to create a zero-knowledge proof. The verifier validates the proof using public parameters and the prover's public values. If the proof is valid, the prover's claim is verified without exposing the secret value.


## Testing

The `cp_zkp_test.go` file contains comprehensive test cases (`TestCPZKPProtocol`) to validate the correctness and soundness of the CP-ZKP protocol. It ensures that a valid proof is correctly verified by the verifier (correctness) and that an invalid proof (soundness) is rejected. 

**TestCPZKPProtocol Function:**
   - `TestCPZKPProtocol` is the main test function for the CP-ZKP protocol.
   - It creates a new instance of `CPZKP`, the protocol handler.
   - It initializes the CP-ZKP parameters (`params1`) using `InitCPZKPParams`.
   - If there's an error during parameter generation, the test will fail, and an error message will be logged.

**Prover Setup:**
   - The prover's secret value `x` is parsed from the configuration file (`sys_config.CPZKP_TEST_X_CORRECT`) as a big integer for testing purposes.
   - A new prover is created using `NewProver(x)`, where `x` is the secret value.

**Prover Creates Proof Commitment:**
   - The prover generates `y1` and `y2` values based on the public parameters and the secret value `x` using `GenerateYValues`.
   - The prover creates the proof commitment `(r1, r2)` using `CreateProofCommitment`.

**Verifier Generates Challenge:**
   - A new verifier is created.
   - The verifier generates a random challenge `c` using `CreateProofChallenge`.
   
**Prover Responds to Challenge and Generates `s`:**
   - The prover creates the response `s` to the verifier's challenge using `CreateProofChallengeResponse`.

**Test Correctness:**
   - The verifier tests the correctness of the proof by calling `VerifyProof`.
   - The verifier checks if `r1` and `r2` match the prover's responses `y1^c * g^s` and `y2^c * h^s`, respectively.
   - If the proof is valid, the test passes; otherwise, it fails with an error message.

**Test Soundness:**
   - An invalid response `invalidS` is created by adding 1 to the valid response `s`.
   - The verifier checks the invalid proof using `VerifyProof`.
   - If the verification returns true (indicating the proof is invalid), the test passes; otherwise, it fails with an error message.

9**TestMain Function:**
   - `TestMain` is responsible for running the tests.
   - The `m.Run()` call executes the tests.

The test code ensures the correctness and soundness of the Chaum-Pedersen Zero-Knowledge Proof (CP-ZKP) protocol. It creates a prover, verifier, and verifies the generated proof against a challenge. The test covers both the correctness (valid proof) and soundness (invalid proof) aspects of the protocol.




