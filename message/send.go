package message

import (
	"context"
	"encoding/json"
	"os"

	"gobox/clients"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

// SendMessage accepts a json.Marshaller item and sends it to the SQS queue url defined by SQS_MESSAGE_QUEUE
func Send(ctx context.Context, clients *clients.Client, item json.Marshaler) (*sqs.SendMessageOutput, error) {
	i, err := item.MarshalJSON()
	mbody := string(i)
	out, err := clients.SQS().SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(os.Getenv("SQS_MESSAGE_QUEUE")),
		MessageBody: &mbody,
	})
	if err != nil {
		logrus.Errorln("error sending sqs message", err)
		return nil, err
	}
	return out, nil
}
