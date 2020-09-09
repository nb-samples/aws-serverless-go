package sample

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// Topic provides SNS client capabilities
type Topic struct {
	Client snsiface.SNSAPI
	ARN    string
}

// Publish an item to the SNS topic and return MessageId
func (t Topic) Publish(item Item) (string, error) {
	// prepare a message body
	body, _ := json.Marshal(item)

	// pubblish the message to SNS topic
	out, err := t.Client.Publish(&sns.PublishInput{
		Message:  aws.String(string(body)),
		Subject:  aws.String("Sample notification message"),
		TopicArn: &t.ARN,
	})

	if err != nil {
		log.Println(err.Error())
		return "", fmt.Errorf("Failed to send a message to %v", t.ARN)
	}
	return *out.MessageId, nil
}

// SnsTopic returns a configured topic client
func SnsTopic(arn string) *Topic {

	sess := session.Must(session.NewSession())

	return &Topic{
		Client: sns.New(sess),
		ARN:    arn,
	}
}
