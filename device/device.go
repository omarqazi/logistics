package device

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

var notificationClient *sns.SNS

func init() {
	config := aws.NewConfig().WithRegion("us-west-2")
	notificationClient = sns.New(session.New(), config)
}

func Push(payload map[string]interface{}, deviceArn string) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	notificationPayload := map[string]string{
		"APNS_SANDBOX": string(payloadBytes),
	}
	finalBytes, err := json.Marshal(notificationPayload)
	if err != nil {
		return err
	}
	params := &sns.PublishInput{
		Message:          aws.String(string(finalBytes)),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(deviceArn),
	}
	_, err = notificationClient.Publish(params)
	if err != nil {
		log.Println("Error sending notification:", err)
		return err
	}

	return nil
}
