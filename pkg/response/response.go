package response

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Response() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Message string `json:"message"`
}

func NewResponse(ctx *fiber.Ctx) *Response {
	return &Response{
		Context: ctx,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	return r
}

func (r *Response) Error(code int, traceId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Message: msg,
	}
	r.IsError = true
	return r
}

func (r *Response) Response() error {
	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes
		}
		return &r.Data
	}())
}
