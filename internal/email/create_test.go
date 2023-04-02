package email

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/harryzcy/mailbox/internal/util/htmlutil"
	"github.com/stretchr/testify/assert"
)

type mockCreateEmailAPI struct {
	mockGetItem            func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	mockPutItem            func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	mockSendEmail          func(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
	mockTransactWriteItems func(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

func (m mockCreateEmailAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.mockGetItem(ctx, params, optFns...)
}

func (m mockCreateEmailAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.mockPutItem(ctx, params, optFns...)
}

func (m mockCreateEmailAPI) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	return m.mockSendEmail(ctx, params, optFns...)
}

func (m mockCreateEmailAPI) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return m.mockTransactWriteItems(ctx, params, optFns...)
}

func TestCreate(t *testing.T) {
	oldGetUpdatedTime := getUpdatedTime
	getUpdatedTime = func() time.Time { return time.Date(2022, 3, 16, 16, 55, 45, 0, time.UTC) }
	defer func() { getUpdatedTime = oldGetUpdatedTime }()

	tableName = "table-for-create"
	tests := []struct {
		client       func(t *testing.T) CreateAndSendEmailAPI
		input        CreateInput
		generateText func(html string) (string, error)
		expected     *CreateResult
		expectedErr  error
	}{
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						t.Helper()

						assert.Equal(t, tableName, *params.TableName)

						messageID := params.Item["MessageID"].(*types.AttributeValueMemberS).Value
						assert.Len(t, messageID, 6+32)
						assert.True(t, strings.HasPrefix(messageID, "draft-"))

						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					Subject: "subject",
					From:    []string{"example@example.com"},
					To:      []string{"example@example.com"},
					Cc:      []string{"example@example.com"},
					Bcc:     []string{"example@example.com"},
					ReplyTo: []string{"example@example.com"},
					Text:    "text",
					HTML:    "<p>html</p>",
				},
				GenerateText: "off",
			},
			expected: &CreateResult{
				TimeIndex: TimeIndex{
					Type:        EmailTypeDraft,
					TimeUpdated: "2022-03-16T16:55:45Z",
				},
				Subject: "subject",
				From:    []string{"example@example.com"},
				To:      []string{"example@example.com"},
				Cc:      []string{"example@example.com"},
				Bcc:     []string{"example@example.com"},
				ReplyTo: []string{"example@example.com"},
				Text:    "text",
				HTML:    "<p>html</p>",
			},
		},
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					Subject: "subject",
					From:    []string{"example@example.com"},
					To:      []string{"example@example.com"},
					Cc:      []string{"example@example.com"},
					Bcc:     []string{"example@example.com"},
					ReplyTo: []string{"example@example.com"},
					HTML:    "<p>example</p>",
				},
				GenerateText: "auto",
			},
			expected: &CreateResult{
				TimeIndex: TimeIndex{
					Type:        EmailTypeDraft,
					TimeUpdated: "2022-03-16T16:55:45Z",
				},
				Subject: "subject",
				From:    []string{"example@example.com"},
				To:      []string{"example@example.com"},
				Cc:      []string{"example@example.com"},
				Bcc:     []string{"example@example.com"},
				ReplyTo: []string{"example@example.com"},
				Text:    "example",
				HTML:    "<p>example</p>",
			},
		},
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					Subject: "subject",
					From:    []string{"example@example.com"},
					To:      []string{"example@example.com"},
					Cc:      []string{"example@example.com"},
					Bcc:     []string{"example@example.com"},
					ReplyTo: []string{"example@example.com"},
					Text:    "text",
					HTML:    "<p>example</p>",
				},
				GenerateText: "auto",
			},
			expected: &CreateResult{
				TimeIndex: TimeIndex{
					Type:        EmailTypeDraft,
					TimeUpdated: "2022-03-16T16:55:45Z",
				},
				Subject: "subject",
				From:    []string{"example@example.com"},
				To:      []string{"example@example.com"},
				Cc:      []string{"example@example.com"},
				Bcc:     []string{"example@example.com"},
				ReplyTo: []string{"example@example.com"},
				Text:    "text",
				HTML:    "<p>example</p>",
			},
		},
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					Subject: "subject",
					From:    []string{"example@example.com"},
					To:      []string{"example@example.com"},
					Cc:      []string{"example@example.com"},
					Bcc:     []string{"example@example.com"},
					ReplyTo: []string{"example@example.com"},
					Text:    "text",
					HTML:    "<p>example</p>",
				},
				GenerateText: "on",
			},
			expected: &CreateResult{
				TimeIndex: TimeIndex{
					MessageID:   "new-message-id",
					Type:        EmailTypeDraft,
					TimeUpdated: "2022-03-16T16:55:45Z",
				},
				Subject: "subject",
				From:    []string{"example@example.com"},
				To:      []string{"example@example.com"},
				Cc:      []string{"example@example.com"},
				Bcc:     []string{"example@example.com"},
				ReplyTo: []string{"example@example.com"},
				Text:    "example",
				HTML:    "<p>example</p>",
			},
		},
		{
			// with Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
					mockSendEmail: func(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						return &sesv2.SendEmailOutput{
							MessageId: aws.String("sent-message-id"),
						}, nil
					},
					mockTransactWriteItems: func(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
						t.Helper()
						assert.Len(t, params.TransactItems, 2)

						for _, item := range params.TransactItems {
							if item.Delete != nil {
								assert.Nil(t, item.Put)
								assert.Equal(t, tableName, *item.Delete.TableName)

								messageID := item.Delete.Key["MessageID"].(*types.AttributeValueMemberS).Value
								assert.Len(t, messageID, 6+32)
								assert.True(t, strings.HasPrefix(messageID, "draft-"))
							}
							if item.Put != nil {
								assert.Nil(t, item.Delete)
								assert.Equal(t, tableName, *item.Put.TableName)

								messageID := item.Put.Item["MessageID"].(*types.AttributeValueMemberS).Value
								assert.Equal(t, "sent-message-id", messageID)
							}
						}

						return &dynamodb.TransactWriteItemsOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					Subject: "subject",
					From:    []string{"example@example.com"},
					To:      []string{"example@example.com"},
					Cc:      []string{"example@example.com"},
					Bcc:     []string{"example@example.com"},
					ReplyTo: []string{"example@example.com"},
					Text:    "text",
					HTML:    "<p>example</p>",
				},
				GenerateText: "on",
				Send:         true,
			},
			expected: &CreateResult{
				TimeIndex: TimeIndex{
					MessageID:   "sent-message-id",
					Type:        EmailTypeSent,
					TimeUpdated: "2022-03-16T16:55:45Z",
				},
				Subject: "subject",
				From:    []string{"example@example.com"},
				To:      []string{"example@example.com"},
				Cc:      []string{"example@example.com"},
				Bcc:     []string{"example@example.com"},
				ReplyTo: []string{"example@example.com"},
				Text:    "example",
				HTML:    "<p>example</p>",
			},
		},
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			input: CreateInput{
				EmailInput:   EmailInput{},
				GenerateText: "on",
			},
			generateText: func(html string) (string, error) {
				return "", errors.New("err")
			},
			expectedErr: errors.New("err"),
		},
		{ // without Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, ErrInvalidInput
					},
				}
			},
			expectedErr: ErrInvalidInput,
		},
		{ // with Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
					mockSendEmail: func(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						return &sesv2.SendEmailOutput{}, errSend
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					From: []string{""},
				},
				Send: true,
			},
			expectedErr: errSend,
		},
		{ // with Send
			client: func(t *testing.T) CreateAndSendEmailAPI {
				return mockCreateEmailAPI{
					mockPutItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return &dynamodb.PutItemOutput{}, nil
					},
					mockSendEmail: func(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
						return &sesv2.SendEmailOutput{MessageId: aws.String("sent-message-id")}, nil
					},
					mockTransactWriteItems: func(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
						return &dynamodb.TransactWriteItemsOutput{}, errBatchWrite
					},
				}
			},
			input: CreateInput{
				EmailInput: EmailInput{
					From: []string{""},
				},
				Send: true,
			},
			expectedErr: errBatchWrite,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()

			if test.generateText != nil {
				generateText = test.generateText
			} else {
				generateText = htmlutil.GenerateText
			}

			actual, err := Create(ctx, test.client(t), test.input)

			if actual != nil && test.expected != nil && !test.input.Send {
				assert.True(t, strings.HasPrefix(actual.MessageID, "draft-"))
				test.expected.MessageID = actual.MessageID // messageID is randomly generated
			}

			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
