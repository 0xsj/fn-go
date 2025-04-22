package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code	int		`json:"code"`
	Message	string	`json:"message"`
	Err		error	`json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:		http.StatusNotFound,
		Message:	message,
		Err:		err,
	}
}

func NewInteral(message string, err error) *AppError {
	return &AppError{
		Code:		http.StatusInternalServerError,
		Message:	message,
		Err:		err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:		http.StatusUnauthorized,
		Message:	message,
		Err:		err,
	}
}