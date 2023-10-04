package contract

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		StatusText: "bad request",
		Message:    err.Error(),
	}
}

func NotFoundErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 404,
		StatusText: "not found",
		Message:    err.Error(),
	}
}

func ServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 500,
		StatusText: "internal server error",
		Message:    "something went wrong, please try again later..",
	}
}

func UnauthorizedErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 401,
		StatusText: "unauthorized",
		Message:    "you are unauthorized to perform this action",
	}
}

func ForbiddenErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 403,
		StatusText: "forbidden",
		Message:    "you are forbidden to perform this action",
	}
}

type RepaymentError struct {
	Msg string
}

func (err RepaymentError) Error() string {
	return err.Msg
}

func NewRepaymentError(msg string) RepaymentError {
	return RepaymentError{Msg: msg}
}
