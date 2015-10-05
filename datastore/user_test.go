package datastore

import (
	"testing"
)

func TestCreateAndFind(t *testing.T) {
	user := User{
		PublicKey: "some-key",
		DeviceId:  "some-id",
		Latitude:  10.0,
		Longitude: 20.0,
	}

	if err := user.Create(); err != nil {
		t.Fatal("Error creating user:", err)
	}

	dbUser, err := GetUser(user.Id)
	if err != nil {
		t.Fatal("Error getting user:", err)
	}

	if dbUser.Id != user.Id || dbUser.PublicKey != user.PublicKey || dbUser.DeviceId != user.DeviceId || dbUser.Latitude != user.Latitude || dbUser.Longitude != user.Longitude {
		t.Fatal("Error: expected user data", user, "but got", dbUser)
	}
}
