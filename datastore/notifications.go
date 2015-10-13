package datastore

func PublishNotification(channel string, payload string) (err error) {
	err = Redis.Publish(channel, payload).Err()
	return
}

func ReceiveMessage(chin string) (chout string, payload string, err error) {
	pubsub, err := Redis.PSubscribe(chin)
	if err != nil {
		return "", "", err
	}
	defer pubsub.Close()

	msg, err := pubsub.ReceiveMessage()
	if err != nil {
		return "", "", err
	}

	return msg.Channel, msg.Payload, nil
}
