package awsapi

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"strconv"
	"sync"
)

func NewSqsAPI() *SqsAPI {

	cfg := &SqsConfig{
		MaxNumbMsg:        int64(10),
		VisibilityTimeout: 30 * 60,
		WaitTimeout:       int64(10),
	}
	return ConfigSqsAPI(cfg)
}

func ConfigSqsAPI(cfg *SqsConfig) *SqsAPI {

	client := sqs.New(LoadAWSConfig())

	return &SqsAPI{
		cfg:    cfg,
		client: client,
	}
}

type SqsConfig struct {
	MaxNumbMsg        int64
	VisibilityTimeout int64
	WaitTimeout       int64
}

type SqsAPI struct {
	cfg    *SqsConfig
	client *sqs.Client
}

func (api *SqsAPI) GetQueueUrl(queueName string) (*string, error) {

	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	req := api.client.GetQueueUrlRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.QueueUrl, nil
}

func (api *SqsAPI) PollMessages(queueUrl string, msgHandler func(string) error) error {

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		WaitTimeSeconds:     aws.Int64(api.cfg.WaitTimeout),
		VisibilityTimeout:   aws.Int64(api.cfg.VisibilityTimeout),
		MaxNumberOfMessages: aws.Int64(api.cfg.MaxNumbMsg),
	}

	req := api.client.ReceiveMessageRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, sqsMsg := range resp.Messages {

		wg.Add(1)
		go func(m *sqs.Message) {
			defer wg.Done()

			err = msgHandler(aws.StringValue(m.Body))
			if err != nil {
				fmt.Printf("handle message failed: %v\n", err)
				return
			}

			err = api.deleteMessage(queueUrl, m)
			if err != nil {
				fmt.Printf("delete message failed: %v\n", err)
			}
		}(&sqsMsg)
	}
	wg.Wait()
	return nil
}

func (api *SqsAPI) ReceiveMessages(queueUrl string) ([]string, error) {

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		WaitTimeSeconds:     aws.Int64(api.cfg.WaitTimeout),
		VisibilityTimeout:   aws.Int64(api.cfg.VisibilityTimeout),
		MaxNumberOfMessages: aws.Int64(api.cfg.MaxNumbMsg),
	}

	req := api.client.ReceiveMessageRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	msgs := []string{}
	entries := []sqs.DeleteMessageBatchRequestEntry{}

	for i, sqsMsg := range resp.Messages {

		msgs[i] = aws.StringValue(sqsMsg.Body)

		id := strconv.Itoa(i)
		entries[i] = sqs.DeleteMessageBatchRequestEntry{
			Id:            &id,
			ReceiptHandle: sqsMsg.ReceiptHandle,
		}
	}

	if len(entries) > 0 {

		go func() {
			delResp, err := api.deleteMessages(queueUrl, entries)
			if err != nil {
				fmt.Printf("delete messages failed: %v\n", err)
				return
			}

			if delResp.Failed != nil && len(delResp.Failed) > 0 {

				for _, entry := range delResp.Failed {

					idx, err := strconv.Atoi(*entry.Id)
					if err != nil {
						fmt.Printf(" convert string to int failed: %v\n", err)
						continue
					}

					respMsg := resp.Messages[idx]
					err = api.deleteMessage(queueUrl, &respMsg)
					if err != nil {
						fmt.Printf("delete message failed: %v\n", err)
						continue
					}
				}
			}
		}()
	}
	return msgs, nil
}

func (api *SqsAPI) deleteMessages(queueUrl string, entries []sqs.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchResponse, error) {

	input := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(queueUrl),
		Entries:  entries,
	}

	req := api.client.DeleteMessageBatchRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *SqsAPI) deleteMessage(queueUrl string, msg *sqs.Message) error {

	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: msg.ReceiptHandle,
	}

	req := api.client.DeleteMessageRequest(input)

	_, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	return nil
}
