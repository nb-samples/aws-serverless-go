package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Proxy constructs a generic API response
func Proxy(r Response) (*events.APIGatewayProxyResponse, error) {
	// new response with a required status
	pxy := events.APIGatewayProxyResponse{StatusCode: r.StatusCode, Headers: make(Headers)}

	// copy headers
	if r.Headers != nil {
		for key, value := range r.Headers {
			pxy.Headers[key] = value
		}
	}

	if r.Body != nil {
		if b, err := json.Marshal(r.Body); err == nil {
			pxy.Headers["Content-Type"] = "application/json"
			pxy.Body = string(b)
		}
	}
	return &pxy, nil
}
