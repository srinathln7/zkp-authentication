package zkp_auth

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ErrInvalidChallengeResponse struct {
	S int64
}

// GRPCStatus : with receiver of type `ErrInvalidCredentials`
// sets the req'd msg using the `status` and `errdetails` pkg
// authentication error `401` is thrown due to invalid credentials
func (e ErrInvalidChallengeResponse) GRPCStatus() *status.Status {
	st := status.New(
		401,
		fmt.Sprintf("authentication error: %d is invalid response to the provided challenge", e.S),
	)

	msg := fmt.Sprintf(
		"Invalid response step provided by the prover (client): %d",
		e.S,
	)

	d := &errdetails.LocalizedMessage{
		Locale:  "en-US",
		Message: msg,
	}

	std, err := st.WithDetails(d)
	if err != nil {
		return st
	}

	return std
}

func (e ErrInvalidChallengeResponse) Error() string {
	return e.GRPCStatus().Err().Error()
}
