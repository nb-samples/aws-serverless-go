package sample

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

// Repo provides DynamoDB client capabilities
type Repo struct {
	Client    dynamodbiface.DynamoDBAPI
	TableName string
}

// Repository returns a configured DynamoDB client
func Repository(tableName string) *Repo {

	sess := session.Must(session.NewSession())

	return &Repo{
		Client:    dynamodb.New(sess),
		TableName: tableName,
	}
}

// Save an item as a new database resource
func (r *Repo) Save(item Item) (*Item, error) {
	if item.ID == "" { // generate a resource id
		item.ID = uuid.New().String()
	}
	now := time.Now()     // set timestamp fields
	item.CreatedAt = &now // reset create timestamp on UPSERT operations
	item.UpdatedAt = &now // reset update timestamp on every change

	// prepare query data
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("Failed to marshal:", err.Error())
		return nil, err
	}
	input := &dynamodb.PutItemInput{Item: av, TableName: &r.TableName}

	// execute query
	if _, err := r.Client.PutItem(input); err != nil {
		log.Println(err.Error())
		return nil, errors.New("Failed to save into the repository")
	}

	return &item, nil
}

// Get an existing resource by ID
func (r *Repo) Get(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, errors.New("Missing resource ID")
	}

	// prepare query data
	input := &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(itemID),
			},
		},
	}

	// execute query
	res, err := r.Client.GetItem(input)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Failed to read item from the repository")
	} else if res.Item == nil {
		return nil, errors.New("Resource not found")
	}

	// process query results
	var item Item
	err = dynamodbattribute.UnmarshalMap(res.Item, &item)
	if err != nil {
		log.Println("Failed to unmarshal:", err.Error())
		return nil, err
	}

	return &item, nil
}

// Delete an existing resource by ID
func (r *Repo) Delete(itemID string) error {
	if itemID == "" {
		return errors.New("Missing resource ID")
	}

	// prepare query data
	input := &dynamodb.DeleteItemInput{
		TableName: &r.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(itemID),
			},
		},
	}

	// execute query
	_, err := r.Client.DeleteItem(input)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to delete item from the repository")
	}

	return nil
}
