package response

import "net/http"

type (
	// Headers map
	Headers map[string]string

	// Response message
	Response struct {
		StatusCode int         `json:"statusCode"` // HTTP status code
		Body       interface{} `json:"body"`       // HTTP response body
		Headers    Headers     `json:"headers"`    // HTTP headers
	}
)

// DefaultStatusText is used for substitution with a standard HTTP status text
const DefaultStatusText string = ""

// Returns status code and text pair.
// An empty message is replaced with a standard HTTP status text.
func httpStatusAs(status int, message string) (int, string) {
	if message == DefaultStatusText {
		message = http.StatusText(status)
	}
	return status, message
}

/**
	SUCCESS RESPONSES
**/

// OK returns 200 status code
func OK(body interface{}, headers Headers) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers:    headers,
	}
}

// Created returns 201 status code
func Created(body interface{}, location string) Response {
	return Response{
		StatusCode: http.StatusCreated,
		Body:       body,
		Headers:    Headers{"Location": location},
	}
}

// NoContent returns 204 status code
func NoContent() Response {
	return Response{
		StatusCode: http.StatusNoContent,
	}
}

/**
	CLIENT ERROR RESPONSES
**/

// BadRequest returns 400 status code
func BadRequest(message string) Response {
	status, message := httpStatusAs(http.StatusBadRequest, message)
	return Response{
		StatusCode: status,
		Body:       Error{Code: status, Message: message},
	}
}

// NotFound returns 404 status code
func NotFound(message string) Response {
	status, message := httpStatusAs(http.StatusNotFound, message)
	return Response{
		StatusCode: status,
		Body:       Error{Code: status, Message: message},
	}
}

// MethodNotAllowed returns 405 status code
func MethodNotAllowed(allow string) Response {
	status, message := httpStatusAs(http.StatusMethodNotAllowed, DefaultStatusText)
	return Response{
		StatusCode: status,
		Body:       Error{Code: status, Message: message},
		Headers:    Headers{"Allow": allow},
	}
}

// Conflict returns 409 status code
func Conflict(message string) Response {
	status, message := httpStatusAs(http.StatusConflict, message)
	return Response{
		StatusCode: status,
		Body:       Error{Code: status, Message: message},
	}
}

/**
	SERVER ERROR RESPONSES
**/

// InternalServerError returns 500 status code
func InternalServerError(message string) Response {
	status, message := httpStatusAs(http.StatusInternalServerError, message)
	return Response{
		StatusCode: status,
		Body:       Error{Code: status, Message: message},
	}
}
