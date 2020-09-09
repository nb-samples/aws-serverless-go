package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nb-samples/aws-serverless-go/response"
	"github.com/stretchr/testify/assert"
)

const (
	uri                = "http://example.com/items"
	invalidJSON        = `"test":"test message"`
	validJSON          = `{"type":"Hello","message":"test message"}`
	validItemWithoutID = `{"name": "unit test", "details": {"description": "test description", "location": "test location", "quantity": 5}}`
	validItemWithID    = `{"id": "test", "name": "unit test", "details": {"description": "test description", "location": "test location", "quantity": 5}}`
)

func TestRouter(t *testing.T) {

	tests := []struct {
		name    string
		request events.APIGatewayProxyRequest
		expect  int
		err     error
	}{
		{
			name: "Negative - DELETE collection",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "DELETE",
				Body:       validJSON,
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
			expect: 405,
		},
		{
			name: "Negative - POST resource",
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     "POST",
				Body:           validJSON,
				Headers:        map[string]string{"Content-Type": "application/json"},
				PathParameters: map[string]string{"itemId": "test-id-value"},
			},
			expect: 405,
		},
		{
			name: "Positive - POST resource",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       validItemWithoutID,
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
			expect: 204,
		},
	}

	var assert = assert.New(t)
	for _, test := range tests {
		res, _ := router(test.request)
		assert.Equal(test.expect, res.StatusCode, "Incorrect status code")
		if res.StatusCode == 405 {
			assert.Contains(res.Headers, "Allow", "Missing HTTP header")
		}
	}
}

func TestCreateFrom(t *testing.T) {

	tests := []struct {
		name    string
		request string
		expect  response.Response
	}{
		{
			name:    "Negative - no content",
			request: "",
			expect:  response.BadRequest(""),
		},
		{
			name:    "Negative - bad request",
			request: invalidJSON,
			expect:  response.BadRequest(""),
		},
		{
			name:    "Negative - has ID",
			request: validItemWithID,
			expect:  response.BadRequest(""),
		},
		{
			name:    "Positive",
			request: validItemWithoutID,
			expect:  response.NoContent(),
		},
	}

	var assert = assert.New(t)
	for _, test := range tests {
		ctx := context.WithValue(context.Background(), keyRequestURI, uri)
		res := createFrom(ctx, test.request)
		assert.Equal(test.expect.StatusCode, res.StatusCode, "Incorrect status code")
		assert.IsType(res.Body, test.expect.Body, "Incorrect body type")
	}
}

func init() {
	isTesting = true
}
