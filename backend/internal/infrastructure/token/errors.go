package token

import (
	"app/pkg/logger"
	"fmt"
)

type ErrorCode int

type ErrorDef struct {
	Code    ErrorCode
	Message string
}

// Error definitions
var (
	ErrTokenGeneration     = ErrorDef{Code: 1020, Message: "error generating token"}
	ErrTokenValidation     = ErrorDef{Code: 1021, Message: "error validating token"}
	ErrTokenMalformed      = ErrorDef{Code: 1022, Message: "malformed token"}
	ErrTokenMissingClaim   = ErrorDef{Code: 1023, Message: "missing required claim"}
	ErrTokenInvalidClaim   = ErrorDef{Code: 1024, Message: "invalid claim value"}
	ErrRefreshTokenExpired = ErrorDef{Code: 1025, Message: "refresh token expired"}
	ErrRefreshTokenInvalid = ErrorDef{Code: 1026, Message: "invalid refresh token"}
	ErrConfiguration       = ErrorDef{Code: 1902, Message: "configuration error"}
)

type ServiceError struct {
	ErrorDef
	Err error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("(%d) %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("(%d) %s", e.Code, e.Message)
}

func NewError(errDef ErrorDef, err error, extra ...string) error {
	message := errDef.Message
	if len(extra) > 0 {
		message = fmt.Sprintf("%s: %s", message, extra[0])
	}
	return &ServiceError{
		ErrorDef: ErrorDef{
			Code:    errDef.Code,
			Message: message,
		},
		Err: err,
	}
}

// Example usage
func LogError(err error) {
	if serviceErr, ok := err.(*ServiceError); ok {
		logger.GetLogger().Error(serviceErr.Error(), map[string]interface{}{
			"code": serviceErr.Code,
		})
	} else {
		logger.GetLogger().Error(err.Error(), map[string]interface{}{})
	}
}
