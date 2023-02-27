package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

// QueryAPI defines set of API required to query for emails
type QueryAPI interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// GetItemAPI defines set of API required to get an email
type GetItemAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// GetItemContentAPI defines set of API required to get attachments or inlines of an email
type GetItemContentAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// DeleteItemAPI defines set of API required to delete an email
type DeleteItemAPI interface {
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// UpdateItemAPI defines set of API required to update an email
type UpdateItemAPI interface {
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

// PutItemAPI defines set of API required to create an new email or replaces an existing email
type PutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// SendEmailAPI defines et of API required to send a email
type SendEmailAPI interface {
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

// SaveAndSendEmailAPI defines set of API required to save an email and send it
type SaveAndSendEmailAPI interface {
	PutItemAPI
	SendEmailAPI
}

// GetAndSendEmailAPI defines set of API required to get and send a email
type GetAndSendEmailAPI interface {
	GetItemAPI
	SendEmailAPI
}

type QueryAndGetItemAPI interface {
	QueryAPI
	GetItemAPI
}

// GetThreadAPI defines set of API required to get a thread and its emails
type GetThreadWithEmailsAPI interface {
	GetItemAPI
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
}

type TransactWriteItemsAPI interface {
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

type StoreEmailAPI interface {
	QueryAPI
	GetItemAPI
	PutItemAPI
	TransactWriteItemsAPI
}
