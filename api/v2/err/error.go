package zkp_auth

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ErrInvalidChallengeResponse struct {
	S string
}

type ErrInvalidRegistration struct {
	User string
}

// GRPCStatus : Sets the req'd msg using the `status` and `errdetails` pkg
// authentication error `401` is thrown due to invalid login credentials
func (e ErrInvalidChallengeResponse) GRPCStatus() *status.Status {

	msg := fmt.Sprintf(
		" ZKP verification failed with the client's challenge response s= : %s",
		e.S,
	)

	st := status.New(
		401,
		"authentication error: invalid login credentials provided",
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

// GRPCStatus : Sets the req'd msg using the `status` and `errdetails` pkg
// registration error `409` is thrown due to duplicate registration of user
func (e ErrInvalidRegistration) GRPCStatus() *status.Status {

	msg := fmt.Sprintf(
		" user %s is already registered on the server",
		e.User,
	)

	st := status.New(
		409,
		"registration error:"+msg,
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

func (e ErrInvalidRegistration) Error() string {
	return e.GRPCStatus().Err().Error()
}
