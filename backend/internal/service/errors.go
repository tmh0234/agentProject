package service

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewError(status int, message string) *AppError {
	return &AppError{Code: status, HTTPStatus: status, Message: message}
}

func WrapError(status int, message string, err error) *AppError {
	return &AppError{Code: status, HTTPStatus: status, Message: message, Err: err}
}

func BadRequest(message string) *AppError {
	return NewError(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return NewError(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return NewError(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return NewError(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return NewError(http.StatusConflict, message)
}

func Validation(message string) *AppError {
	return NewError(http.StatusBadRequest, message)
}

func Internal(message string, err error) *AppError {
	return WrapError(http.StatusInternalServerError, message, err)
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func IsUnauthorized(err error) bool {
	appErr, ok := AsAppError(err)
	return ok && appErr.HTTPStatus == http.StatusUnauthorized
}

func IsConflict(err error) bool {
	appErr, ok := AsAppError(err)
	return ok && appErr.HTTPStatus == http.StatusConflict
}

func IsValidation(err error) bool {
	appErr, ok := AsAppError(err)
	return ok && appErr.HTTPStatus == http.StatusBadRequest
}
