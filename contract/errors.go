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

var (
	ErrNotFound   = &ErrorResponse{StatusCode: 404, Message: "not found"}
	ErrBadRequest = &ErrorResponse{StatusCode: 400, Message: "bad request"}
)

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
		Message:    err.Error(),
	}
}
