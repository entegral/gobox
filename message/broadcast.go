package message

import (
	"context"
	"encoding/json"
	"os"

	"gobox/clients"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/sirupsen/logrus"
)

// BroadcastMessage accepts a json.Marshaller item and a title string. The item marshals itself into byte data and attaches
// itself as a string to the "detail" portion of an EventBridge event. EventBridge can trigger various other aws services
// using pattern matching against the fields of this "detail" object, as well as against the title provided as an argument
// to this function. It automatically broadcasts these events to the bus defined by the EVENT_BUS_NAME env var and it also
// attaches the AWS_LAMBDA_FUNCTION_NAME as the source of the event. Ensure your EB rules exist on the EVENT_BUS_NAME bus.
func BroadcastMessage(ctx context.Context, clients *clients.Client, item json.Marshaler, title string) (*eventbridge.PutEventsOutput, error) {
	detailBytes, err := item.MarshalJSON()
	if err != nil {
		logrus.Errorln("error marshalling item in BroadcastMessage", err)
		return nil, err
	}
	out, err := clients.EventBridge().PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:       aws.String(string(detailBytes)),
				DetailType:   aws.String(title),
				Source:       aws.String(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")),
				EventBusName: aws.String(os.Getenv("EVENT_BUS_NAME")),
			},
		},
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"EventBus": os.Getenv("EVENT_BUS_NAME"),
		}).Errorln("error putting event onto event bus", err)
		return nil, err
	}
	return out, nil
}
