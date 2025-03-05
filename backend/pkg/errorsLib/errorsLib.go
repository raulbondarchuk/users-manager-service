package errorsLib

import (
	"errors"
	"net/http"
)

var (
	ErrAccessDenied = errors.New("access denied")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
)

// HTTPStatusCode returns the HTTP status code based on the error message
func HTTPStatusCode(err string) int {
	switch err {
	case ErrAccessDenied.Error():
		return http.StatusUnauthorized // 401
	case ErrForbidden.Error():
		return http.StatusForbidden // 403
	case ErrNotFound.Error():
		return http.StatusNotFound // 404
	default:
		return http.StatusInternalServerError // 500
	}
}
