package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code          int    `json:"code"`
	Err           string `json:"error"`
	ValidateError any    `json:"validateError,omitempty"`
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func NewValidationError(errors any) Error {
	return Error{
		Code:          http.StatusBadRequest,
		Err:           "invalid payload",
		ValidateError: errors,
	}
}

func (e Error) Error() string {
	return e.Err
}

func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "not a valid id")
}

func ErrUnauthorized(m ...string) Error {
	var msg = "unauthorized request"
	if len(m) != 0 {
		msg = m[0]
	}
	return NewError(http.StatusUnauthorized, msg)
}

func ErrBadJSON() Error {
	return NewError(http.StatusBadRequest, "bad json")
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if err, ok := err.(Error); ok {
		return ctx.Status(err.Code).JSON(err)
	}

	internalError := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(internalError.Code).JSON(internalError)
}
