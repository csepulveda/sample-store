package sqs

import (
	"context"
	"encoding/json"
	"log"
	"products-worker/internal/processor"
	"products-worker/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SNSMessageWrapper struct {
	Message           string                         `json:"Message"`
	MessageAttributes map[string]SNSMessageAttribute `json:"MessageAttributes"`
}

type SNSMessageAttribute struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

func ListenAndProcess(ctx context.Context, client *sqs.Client, queueURL string, handler processor.Handler) {
	for {
		output, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("error receiving messages: %v", err)
			continue
		}

		for _, msg := range output.Messages {
			if err := processMessage(ctx, client, queueURL, msg, handler); err != nil {
				log.Printf("processing failed: %v", err)
			}
		}
	}
}

func processMessage(ctx context.Context, client *sqs.Client, queueURL string, msg types.Message, handler processor.Handler) error {
	var sns SNSMessageWrapper
	if err := json.Unmarshal([]byte(*msg.Body), &sns); err != nil {
		return err
	}

	traceparent, _ := sns.MessageAttributes["traceparent"]
	ctx, span := tracing.NewSpanWithTraceparent(ctx, "processMessage", traceparent.Value)
	defer span.End()

	if err := handler.HandleMessage(ctx, sns.Message); err != nil {
		return err
	}

	_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}
