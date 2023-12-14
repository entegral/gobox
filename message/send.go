package message

import (
	"context"
	"encoding/json"

	"github.com/entegral/gobox/clients"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

// SendMessage accepts an item and sends it to the SQS queue url provided by the queueURL argument
func Send(ctx context.Context, queueURL string, item any) (*sqs.SendMessageOutput, error) {
	client := clients.GetDefaultClient(ctx)
	return SendWithClient(ctx, client, queueURL, item)
}

// SendWithClient accepts any item and sends it to the SQS queue url provided by the queueURL argument
// using the client provided by the ctx argument
func SendWithClient(ctx context.Context, client *clients.Client, queueURL string, item any) (*sqs.SendMessageOutput, error) {
	i, err := json.Marshal(item)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Item": item,
		}).Errorln("error marshalling item into json", err)
		return nil, err
	}
	mbody := string(i)
	out, err := client.SQS().SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &mbody,
	})
	if err != nil {
		logrus.Errorln("error sending sqs message", err)
		return nil, err
	}
	return out, nil
}
