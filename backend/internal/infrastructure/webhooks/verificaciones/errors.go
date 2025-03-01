package verificaciones

import (
	"app/pkg/logger"
	"fmt"
)

type ErrorCode int

type ErrorDef struct {
	Code    ErrorCode
	Message string
}

// Ваши ErrorDef:
var (
	// Static errors
	ErrRequestFailed           = ErrorDef{Code: 3000, Message: "request failed"}
	ErrLoginFailed             = ErrorDef{Code: 3001, Message: "verificaciones user login failed"}
	ErrInvalidLoginResponse    = ErrorDef{Code: 3002, Message: "verificaciones user login failed: invalid response"}
	ErrGetCompanyFailed        = ErrorDef{Code: 3003, Message: "failed to get company by company id"}
	ErrTokenFailed             = ErrorDef{Code: 3004, Message: "failed to get token"}
	ErrParseResponse           = ErrorDef{Code: 3005, Message: "failed to parse response"}
	ErrCheckUserExistsFailed   = ErrorDef{Code: 3006, Message: "check if user exists request failed"}
	ErrCompanyRequestFailed    = ErrorDef{Code: 3007, Message: "company request failed"}
	ErrParseCompanyResponse    = ErrorDef{Code: 3008, Message: "failed to parse company response"}
	ErrParseUserExistsResponse = ErrorDef{Code: 3009, Message: "failed to parse check if user exists response"}
	ErrEmptyTokenResponse      = ErrorDef{Code: 3010, Message: "empty token in response"}
	ErrLoginRequestFailed      = ErrorDef{Code: 3011, Message: "login request failed"}
	ErrParseLoginResponse      = ErrorDef{Code: 3012, Message: "failed to parse login response"}
	ErrICCIDRequestFailed      = ErrorDef{Code: 3013, Message: "ICCID request failed"}
	ErrParseICCIDResponse      = ErrorDef{Code: 3014, Message: "failed to parse ICCID response"}
	ErrSecurityFailed          = ErrorDef{Code: 3015, Message: "access denied"}

	// Dynamic errors
	ErrUnexpectedResponse   = ErrorDef{Code: 3200, Message: "unexpected response"}
	ErrExternalServiceError = ErrorDef{Code: 3201, Message: "external service error"}
)

type ServiceError struct {
	ErrorDef
	Err error
}

// If you have a logging system, for example logger.GetLogger(),
// then inside Error() you can log.
func (e *ServiceError) Error() string {
	if e.Err != nil {
		logger.GetLogger().VerificacionesError(e.Message, map[string]interface{}{"error": e.Err})
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	logger.GetLogger().VerificacionesWarn("Verificaciones error", map[string]interface{}{"error": e.Message, "code": e.Code})
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewServiceError — convenient helper for creating errors
func NewServiceError(errDef ErrorDef, err error) error {
	return &ServiceError{
		ErrorDef: errDef,
		Err:      err,
	}
}
