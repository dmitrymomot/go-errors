package errors

import (
	"fmt"
	"net/http"
	"net/url"
)

type (
	// Error struct
	err struct {
		Code   int         `json:"code"`
		Title  string      `json:"title"`
		Detail interface{} `json:"detail,omitempty"`
	}
)

func (e err) Error() string {
	if e.Detail != nil {
		return fmt.Sprintf("%+v", e.Detail)
	} else if e.Title != "" {
		return e.Title
	}
	if e.Code <= 0 {
		e.Code = http.StatusInternalServerError
	}
	return http.StatusText(e.Code)
}

// New error
func New(e interface{}) error {
	if er, ok := e.(error); ok {
		return newError(0, "", er.Error())
	}
	return newError(0, "", e)
}

// NewValidation error
func NewValidation(e url.Values) error {
	return newError(http.StatusUnprocessableEntity, "Validation Error", e)
}

// NewHTTP error
func NewHTTP(code int, e interface{}) error {
	return newError(code, "", e)
}

// WrapHTTP wraps error
func WrapHTTP(e error) error {
	if er, ok := e.(err); ok {
		return er
	}
	if er, ok := e.(error); ok {
		return newError(0, "", er.Error())
	}
	return newError(0, "", e)
}

func newError(code int, title string, detail interface{}) err {
	if code < 400 || code > 599 {
		code = http.StatusInternalServerError
	}
	if title == "" {
		title = http.StatusText(code)
	}
	return err{
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}
