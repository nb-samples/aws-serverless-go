package sample

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type mockDdb struct {
	dynamodbiface.DynamoDBAPI
	item *Item
	err  error
}

func (mock *mockDdb) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, mock.err
}

func (mock *mockDdb) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	output := new(dynamodb.GetItemOutput)
	if mock.item != nil {
		output.Item, _ = dynamodbattribute.MarshalMap(&mock.item)
	}
	return output, mock.err
}

func (mock *mockDdb) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return nil, mock.err
}

func TestRepo_Save(t *testing.T) {
	type fields struct {
		Client    dynamodbiface.DynamoDBAPI
		TableName string
	}
	type args struct {
		item Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Item
		wantErr bool
	}{
		{
			name: "successful operation with UUID",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args: args{item: Item{Name: "test-item-name"}},
			want: &Item{Name: "test-item-name"},
		},
		{
			name: "successful operation with existing ID",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args: args{item: Item{ID: "test-item-id", Name: "test-item-name"}},
			want: &Item{ID: "test-item-id", Name: "test-item-name"},
		},
		{
			name: "failed operation",
			fields: fields{
				Client:    &mockDdb{err: errors.New("Mock DynamoDB error")},
				TableName: "mock-table",
			},
			args:    args{item: Item{Name: "test-item-name"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				Client:    tt.fields.Client,
				TableName: tt.fields.TableName,
			}

			got, err := r.Save(tt.args.item)

			assert := assert.New(t)
			if tt.wantErr {
				assert.Error(err)

			} else if assert.NoError(err) {
				assert.NotSame(tt.want, got, "Returned original item")
				assert.NotEmpty(got.ID, "UUID")
				if assert.NotNil(got.CreatedAt, "CreatedAt") &&
					assert.NotNil(got.UpdatedAt, "UpdatedAt") {
					assert.EqualValues(got.CreatedAt, got.UpdatedAt, "CreatedAt, UpdatedAt")
				}

				assert.Equal(tt.want.Name, got.Name, "Name")
				if tt.want.ID != "" {
					assert.Equal(tt.want.ID, got.ID, "ID")
				}
			}
		})
	}
}

func TestRepo_Get(t *testing.T) {
	type fields struct {
		Client    dynamodbiface.DynamoDBAPI
		TableName string
	}
	type args struct {
		itemID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Item
		wantErr bool
	}{
		{
			name: "successful operation",
			fields: fields{
				Client:    &mockDdb{item: &Item{ID: "test-item-id", Name: "test-item-name"}},
				TableName: "mock-table",
			},
			args: args{itemID: "test-item-id"},
			want: &Item{ID: "test-item-id", Name: "test-item-name"},
		},
		{
			name: "item not found",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args:    args{itemID: "test-item-id"},
			wantErr: true,
		},
		{
			name: "bad request",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args:    args{itemID: ""},
			wantErr: true,
		},
		{
			name: "failed operation",
			fields: fields{
				Client:    &mockDdb{err: errors.New("Mock DynamoDB error")},
				TableName: "mock-table",
			},
			args:    args{itemID: "test-item-id"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				Client:    tt.fields.Client,
				TableName: tt.fields.TableName,
			}

			got, err := r.Get(tt.args.itemID)

			assert := assert.New(t)
			if tt.wantErr {
				assert.Error(err)

			} else if assert.NoError(err) {
				assert.Equal(tt.want.ID, got.ID, "ID")
				assert.Equal(tt.want.Name, got.Name, "Name")
			}
		})
	}
}

func TestRepo_Delete(t *testing.T) {
	type fields struct {
		Client    dynamodbiface.DynamoDBAPI
		TableName string
	}
	type args struct {
		itemID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful operation",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args:    args{itemID: "test-item-id"},
			wantErr: false,
		},
		{
			name: "bad request",
			fields: fields{
				Client:    &mockDdb{},
				TableName: "mock-table",
			},
			args:    args{itemID: ""},
			wantErr: true,
		},
		{
			name: "failed operation",
			fields: fields{
				Client:    &mockDdb{err: errors.New("Mock DynamoDB error")},
				TableName: "mock-table",
			},
			args:    args{itemID: "test-item-id"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				Client:    tt.fields.Client,
				TableName: tt.fields.TableName,
			}

			err := r.Delete(tt.args.itemID)

			assert := assert.New(t)
			if tt.wantErr {
				assert.Error(err)

			} else {
				assert.NoError(err)
			}
		})
	}
}
