package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {

	tests := []struct {
		name       string
		statusCode int
		headers    map[string]string
		body       interface{}
		jsonBody   string
	}{
		{
			name:       "HTTP status, headers and JSON body",
			statusCode: 200,
			headers:    map[string]string{"x-test": "test header"},
			body:       struct{ Message string }{"test message"},
			jsonBody:   `{"Message":"test message"}`,
		},
		{
			name:       "without body",
			statusCode: 204,
			headers:    map[string]string{"x-test": "test header"},
		},
		{
			name:       "HTTP status only",
			statusCode: 204,
		},
	}

	assert := assert.New(t)

	for _, test := range tests {
		pxy, _ := Proxy(Response{StatusCode: test.statusCode, Body: test.body, Headers: test.headers})

		assert.Equal(test.statusCode, pxy.StatusCode, "Incorrect status code")

		for k, v := range test.headers {
			assert.Contains(pxy.Headers, k, "Missing HTTP header: %v", k)
			assert.Equal(v, pxy.Headers[k], "Incorrect HTTP header: %v", k)
		}

		if test.body == nil {

			assert.NotContains(pxy.Headers, "Content-Type", "Unexpected header: %v", "Content-Type")

		} else if test.jsonBody != "" {

			assert.Equal(test.jsonBody, pxy.Body, "Incorrect response body")

			for k, v := range map[string]string{"Content-Type": "application/json"} {
				assert.Contains(pxy.Headers, k, "Missing HTTP header: %v", k)
				assert.Equal(v, pxy.Headers[k], "Incorrect HTTP header: %v", k)
			}
		}
	}
}
