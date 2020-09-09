package sample

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/stretchr/testify/assert"
)

// Mock SNS client with successful response
type mockSns struct {
	snsiface.SNSAPI
	msgID string
	err   error
}

func (mock *mockSns) Publish(*sns.PublishInput) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{MessageId: &mock.msgID}, mock.err
}

func TestTopic_Publish(t *testing.T) {
	type fields struct {
		Client snsiface.SNSAPI
		ARN    string
	}
	type args struct {
		item Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "successful publishing",
			fields: fields{
				Client: &mockSns{msgID: "test-message-id"},
				ARN:    "arn:mock:sns:topic",
			},
			args: args{item: Item{}},
			want: "test-message-id",
		},
		{
			name: "unsuccessful publishing",
			fields: fields{
				Client: &mockSns{err: errors.New("Mock SNS error")},
				ARN:    "arn:mock:sns:topic",
			},
			args:    args{item: Item{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic := Topic{
				Client: tt.fields.Client,
				ARN:    tt.fields.ARN,
			}

			got, err := topic.Publish(tt.args.item)

			assert := assert.New(t)
			if tt.wantErr {
				assert.Error(err)

			} else if assert.NoError(err) {
				assert.Equal(tt.want, got, "MessageId")
			}
		})
	}
}
