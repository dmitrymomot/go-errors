package errors

import (
	"fmt"
	"net/http"
	"net/url"
)

type (
	// Error interface
	Error interface {
		error
		GetCode() int
		GetTitle() string
		GetDetail() interface{}
	}
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

func (e err) GetCode() int {
	return e.Code
}

func (e err) GetTitle() string {
	return e.Title
}

func (e err) GetDetail() interface{} {
	return e.Detail
}

// New error
func New(e interface{}) Error {
	if er, ok := e.(error); ok {
		return newError(0, "", er.Error())
	}
	return newError(0, "", e)
}

// NewValidation error
func NewValidation(e url.Values) Error {
	return newError(http.StatusUnprocessableEntity, "Validation Error", e)
}

// NewHTTP error
func NewHTTP(code int, e interface{}) Error {
	return newError(code, "", e)
}

// WrapHTTP wraps error
func WrapHTTP(e error) Error {
	if er, ok := e.(Error); ok {
		return er
	}
	if er, ok := e.(error); ok {
		return newError(0, "", er.Error())
	}
	return newError(0, "", e)
}

// Wrap error
func Wrap(e error, message string) error {
	if e == nil {
		return nil
	}
	if er, ok := e.(Error); ok {
		return newError(er.GetCode(), er.GetTitle(), fmt.Sprintf("%s: %s", message, er.Error()))
	}
	return newError(0, "", fmt.Sprintf("%s: %s", message, e.Error()))
}

func newError(code int, title string, detail interface{}) Error {
	if er, ok := detail.(Error); ok {
		return er
	}

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
