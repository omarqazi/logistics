package datastore

import (
	"fmt"
	"testing"
	"time"
)

func TestNotifications(t *testing.T) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		PublishNotification("test-channel", "a-value")
	}()

	timeout := time.Tick(200 * time.Millisecond)
	finishChannel := make(chan error, 1)

	go func() {
		x, y, err := ReceiveMessage("test-channel")
		if x != "test-channel" || y != "a-value" {
			finishChannel <- fmt.Errorf("Expected test-channle a-value but got %s %s", x, y)
			return
		}

		if err != nil {
			finishChannel <- err
		} else {
			finishChannel <- nil
		}
	}()

	select {
	case <-timeout:
		t.Error("Error: timed out before notification received")
	case err := <-finishChannel:
		if err != nil {
			t.Fatal("Error receiving message:", err)
		}
	}
}
