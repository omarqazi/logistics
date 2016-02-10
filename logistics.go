package main

import (
	"fmt"
	"github.com/omarqazi/logistics/device"
	"net/http"
)

func sendPushNotification() {
	pushPayload := map[string]interface{}{
		"aps": map[string]interface{}{
			"sound": "supreme.caf",
			"alert": map[string]string{
				"title": "New Order",
				"body":  "Hannah wants to rent Superbad",
			},
			"category": "newOrder",
		},
	}
	endpoint := "arn:aws:sns:us-west-2:716704604621:endpoint/APNS_SANDBOX/Market/25dd5687-565a-33f7-ab3e-d42c5c066f43"
	device.Push(pushPayload, endpoint)
}

func main() {
	sendPushNotification()
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
