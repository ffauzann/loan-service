package logging

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Code    codes.Code
	Message string
}

// Error implements required method to fulfill error interface.
func (e *errorResponse) Error() string {
	return e.Message
}

// GRPCStatus implements required method to fulfill anonym interface.
// Used in status.FromError().
func (e *errorResponse) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}
