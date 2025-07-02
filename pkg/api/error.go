package api

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	Code          int
	ResponseCode  int
	Message       string
	InternalError error
}

func (s ServiceError) Error() string {
	if s.InternalError != nil {
		return fmt.Sprintf(`{"error": "%d - %s - Service Response code: %d", "internalError": "%s"}`, s.Code, s.Message, s.ResponseCode, s.InternalError.Error())
	}
	return fmt.Sprintf(`{"error": "%d - %s - Service Response code: %d"}`, s.Code, s.Message, s.ResponseCode)
}

const (
	ErrNotFound = 1 + iota
	ErrNotImplemented
	ErrBadRequest
	ErrInternalServer
	ErrUnatuhorize
	ErrConflict
	ErrUnprocessableEntity
)

var ServiceErrors = map[int]ServiceError{
	ErrNotFound: {
		Code:         ErrNotFound,
		ResponseCode: http.StatusNotFound,
		Message:      "document not found",
	},
	ErrNotImplemented: {
		Code:         ErrNotImplemented,
		ResponseCode: http.StatusUnavailableForLegalReasons,
		Message:      "not implemented yet!",
	},
	ErrBadRequest: {
		Code:         ErrBadRequest,
		ResponseCode: http.StatusBadRequest,
		Message:      "bad request",
	},
	ErrInternalServer: {
		Code:         ErrInternalServer,
		ResponseCode: http.StatusInternalServerError,
		Message:      "Internal Server Error",
	},
	ErrUnatuhorize: {
		Code:         ErrUnatuhorize,
		ResponseCode: http.StatusUnauthorized,
		Message:      "unauthorize token",
	},
	ErrConflict: {
		Code:         ErrConflict,
		ResponseCode: http.StatusConflict,
		Message:      "conflict",
	},
	ErrUnprocessableEntity: {
		Code:         ErrUnprocessableEntity,
		ResponseCode: http.StatusUnprocessableEntity,
		Message:      "unprocessable Entity",
	},
}
