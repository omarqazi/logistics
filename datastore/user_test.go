package datastore

import (
	"testing"
)

func TestCRUD(t *testing.T) {
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

	user.PublicKey = "changed-key"
	user.Latitude = 50.5
	user.Longitude = 60.6
	if err := user.Update(); err != nil {
		t.Fatal("Error updating user:", err)
	}

	dbUser, err = GetUser(user.Id)
	if dbUser.Id != user.Id || dbUser.PublicKey != user.PublicKey || dbUser.DeviceId != user.DeviceId || dbUser.Latitude != user.Latitude || dbUser.Longitude != user.Longitude {
		t.Fatal("Error: expected updated user data", user, "but got", dbUser)
	}

	if err := user.Delete(); err != nil {
		t.Fatal("Error deleting user:", err)
	}

	dbUser, err = GetUser(user.Id)
	if err != nil {
		t.Fatal("Error getting user after delete:", err)
	}
	if dbUser != nil {
		t.Fatal("Error: tried deleting user but got", dbUser)
	}
}
