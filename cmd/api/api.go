package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nb-samples/aws-serverless-go/internal/sample"
	"github.com/nb-samples/aws-serverless-go/response"
)

const (
	envTableName = "DB_TABLE_NAME"
	envTopicArn  = "SNS_TOPIC_ARN"
)

type configuration struct {
	dbTableName string
	snsTopicArn string
}

func (c *configuration) incomplete() bool {
	return c.dbTableName == "" || c.snsTopicArn == ""
}

var (
	isTesting bool
	config    configuration
)

type key int

const (
	keyRequestURI key = iota + 1
)

// Request router
func router(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var resp response.Response
	ctx := context.WithValue(context.Background(), keyRequestURI, requestURI(req))

	itemID := req.PathParameters["itemId"]
	if itemID == "" {
		// collection actions
		switch req.HTTPMethod {
		case "POST":
			resp = createFrom(ctx, req.Body)
		default:
			resp = response.MethodNotAllowed("POST")
		}
	} else {
		// resource actions
		switch req.HTTPMethod {
		case "GET":
			resp = get(ctx, itemID)
		case "DELETE":
			resp = delete(ctx, itemID)
		default:
			resp = response.MethodNotAllowed("GET, DELETE")
		}
	}
	log.Println("RESPONSE:", resp)
	return response.Proxy(resp)
}

func requestURI(req events.APIGatewayProxyRequest) string {
	proto := req.Headers["X-Forwarded-Proto"]
	host := req.Headers["Host"]
	path := req.Path
	if strings.Contains(host, ".execute-api.") {
		stage := req.RequestContext.Stage
		return proto + "://" + host + "/" + stage + path
	}
	return proto + "://" + host + path
}

// Creates a new resource. ID is auto allocated and not allowed in the message.
func createFrom(ctx context.Context, body string) response.Response {
	var item sample.Item

	if err := json.Unmarshal([]byte(body), &item); err != nil {
		log.Println(err.Error())
		return response.BadRequest(err.Error())
	}

	if item.ID != "" {
		return response.BadRequest("Item ID is not allowed when creating a new resource.")
	}

	if isTesting {
		// FIXME unit tests exit here
		return response.NoContent()

	} else if config.incomplete() {
		log.Fatalln("Service is not configured")
	}

	// publish item to SNS topic
	if msgID, err := sample.SnsTopic(config.snsTopicArn).Publish(item); err == nil {
		fmt.Println("SNS notification:", msgID)
	}

	// save item in DynamoDB
	out, err := sample.Repository(config.dbTableName).Save(item)
	if err != nil {
		// return 400 Bad Request for simplicity
		return response.BadRequest(err.Error())
	}
	fmt.Println("DynamoDB persistence:", out.ID)
	return response.Created(out, ctx.Value(keyRequestURI).(string)+"/"+out.ID)
}

// Gets a resource by ID
func get(ctx context.Context, itemID string) response.Response {
	if isTesting {
		// unit tests exit here
		return response.NoContent()
	}

	out, err := sample.Repository(config.dbTableName).Get(itemID)
	if err != nil {
		// return 404 Not Found for simplicity
		return response.NotFound(err.Error())
	}
	return response.OK(out, nil)
}

// Deletes a resource by ID.
func delete(ctx context.Context, itemID string) response.Response {
	if isTesting {
		// unit tests exit here
		return response.NoContent()
	}

	if err := sample.Repository(config.dbTableName).Delete(itemID); err != nil {
		// return 404 Not Found for simplicity
		return response.NotFound(err.Error())
	}
	return response.NoContent()
}

func init() {
	var ok bool

	if config.dbTableName, ok = os.LookupEnv(envTableName); !ok {
		log.Println("Missing environment variable:", envTableName)
	}
	if config.snsTopicArn, ok = os.LookupEnv(envTopicArn); !ok {
		log.Println("Missing environment variable:", envTopicArn)
	}
}

func main() {
	// Make the handler available for RPC by AWS Lambda
	lambda.Start(router)
}
