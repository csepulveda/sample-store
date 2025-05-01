package publisher

import (
	"context"
	"encoding/json"
	"log"
	"orders-service/internal/domain"
	"orders-service/internal/tracing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SnsOrderPublisher struct {
	client   *sns.Client
	topicArn string
}

func NewSnsOrderPublisher(client *sns.Client, topicArn string) *SnsOrderPublisher {
	return &SnsOrderPublisher{
		client:   client,
		topicArn: topicArn,
	}
}

func (p *SnsOrderPublisher) PublishOrderCreated(ctx context.Context, order domain.Order) error {
	return p.publish(ctx, "order.created", order)
}

func (p *SnsOrderPublisher) PublishOrderCanceled(ctx context.Context, order domain.Order) error {
	return p.publish(ctx, "order.canceled", order)
}

func (p *SnsOrderPublisher) publish(ctx context.Context, eventType string, order domain.Order) error {
	ctx, span := tracing.NewSpan(ctx, "SnsOrderPublisher#publish")
	defer span.End()

	span.SetAttributes(
		tracing.StringAttribute("orderId", order.ID),
		tracing.StringAttribute("eventType", eventType),
	)
	traceparent := tracing.GetTraceParent(ctx)

	messageAttributes := map[string]types.MessageAttributeValue{
		"traceparent": {
			DataType:    aws.String("String"),
			StringValue: aws.String(traceparent),
		},
	}

	payload := map[string]interface{}{
		"type":     eventType,
		"orderId":  order.ID,
		"items":    order.Items,
		"datetime": order.CreatedAt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		TopicArn:          aws.String(p.topicArn),
		Message:           aws.String(string(body)),
		MessageAttributes: messageAttributes,
	})

	log.Printf("Publishing attribute: %s", traceparent)
	return err
}
