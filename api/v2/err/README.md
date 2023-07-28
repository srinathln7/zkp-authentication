## Custom GRPC Errors in  Package `zkp_auth`:

The `zkp_auth` package defines custom grpc error types and their corresponding `GRPCStatus` functions to handle specific authentication errors related to ZKP Authentication. These error types provide more specific information about the nature of authentication errors, such as invalid challenge responses or duplicate user registrations. The `GRPCStatus` function generate GRPC `status.Status` instances with detailed error messages that can be used to convey meaningful error information to clients and consumers of the ZKP Authentication service.

1. **ErrInvalidChallengeResponse:**
   - This is a custom error type that represents an invalid challenge response during ZKP verification.
   - It contains a field `S` representing the client's challenge response `s`.

2. **ErrInvalidRegistration:**
   - This is a custom error type that represents an invalid user registration attempt.
   - It contains a field `User` representing the username for which registration failed.

3. **GRPCStatus for ErrInvalidChallengeResponse:**
   - The `GRPCStatus()` method of `ErrInvalidChallengeResponse` returns a GRPC `status.Status` instance with specific error details.
   - It sets the error code to `401`, indicating an authentication error due to invalid login credentials.
   - It constructs a detailed error message using the `fmt.Sprintf` function, including the client's challenge response `s`.
   - The method creates a `errdetails.LocalizedMessage` object containing the detailed error message.
   - It attaches the `errdetails.LocalizedMessage` to the `status.Status` using the `WithDetails` function.
   - If there is an error while attaching details, the method returns the original `status.Status`.

4. **Error for ErrInvalidChallengeResponse:**
   - The `Error()` method of `ErrInvalidChallengeResponse` returns the error message from the corresponding GRPCStatus using `e.GRPCStatus().Err().Error()`.

5. **GRPCStatus for ErrInvalidRegistration:**
   - The `GRPCStatus()` method of `ErrInvalidRegistration` returns a GRPC `status.Status` instance with specific error details.
   - It sets the error code to `409`, indicating a registration error due to duplicate registration of a user.
   - It constructs a detailed error message using the `fmt.Sprintf` function, including the username that caused the registration failure.
   - The method creates a `errdetails.LocalizedMessage` object containing the detailed error message.
   - It attaches the `errdetails.LocalizedMessage` to the `status.Status` using the `WithDetails` function.
   - If there is an error while attaching details, the method returns the original `status.Status`.

6. **Error for ErrInvalidRegistration:**
   - The `Error()` method of `ErrInvalidRegistration` returns the error message from the corresponding GRPCStatus using `e.GRPCStatus().Err().Error()`.




