package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/omarqazi/logistics/datastore"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUser(t *testing.T) {
	resp := httptest.NewRecorder()
	uri := "/users/whatever"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 500 {
		t.Fatal("Expected 500 error requesting invalid UUID but got", resp.Code)
	}

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else {
		if !strings.Contains(string(p), "Error") {
			t.Error("User controller error response does not contain word 'Error'")
		}
	}

	coolUser := datastore.User{PublicKey: "something"}
	if err := coolUser.Create(); err != nil {
		t.Fatal("Error creating user:", err)
	}

	resp = httptest.NewRecorder()
	uri = fmt.Sprintf("/users/%s", coolUser.Id)
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatal("Expected 200 but got", resp.Code)
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	dbUser := datastore.User{}
	if err := json.Unmarshal(p, &dbUser); err != nil {
		t.Fatal("Error decoding response json", err)
	}

	expectedContentType := "application/json"
	if ct := resp.Header().Get("Content-Type"); ct != expectedContentType {
		t.Error("Expected content type", expectedContentType, "but got", ct)
	}

	if dbUser.Id != coolUser.Id || dbUser.PublicKey != coolUser.PublicKey {
		t.Error("Error: expected user", coolUser, "but got", dbUser)
	}
}

func TestGetMissingUser(t *testing.T) {
	resp := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", datastore.NewUUID())
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 404 {
		t.Fatal("Expected 404 but got", resp.Code)
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(p), "does not exist") {
		t.Error("User controller 404 response does not contain phrase 'does not exist'", string(p))
	}
}

func TestPostUser(t *testing.T) {
	resp := httptest.NewRecorder()
	uri := "/users/"
	user := datastore.User{
		PublicKey: "something",
		Latitude:  10.0,
		Longitude: 11.0,
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal("Error marshaling json:", err)
	}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatal("Expected 200 on user create but got", resp.Code)
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var parsedUser datastore.User
	if err := json.Unmarshal(p, &parsedUser); err != nil {
		t.Fatal(err)
	}

	if parsedUser.Id == "" {
		t.Error("Error: parsed user has no ID")
	}

	if parsedUser.PublicKey != user.PublicKey {
		t.Error("Error: expected public key", user.PublicKey, "but got", parsedUser.PublicKey)
	}

	dbUser, err := datastore.GetUser(parsedUser.Id)
	if err != nil {
		t.Error("Error getting posted user:", err)
		return
	}
	if dbUser == nil {
		t.Error("Posted user not found in database")
	}

	if err := dbUser.Delete(); err != nil {
		t.Error("Error deleting user:", err)
	}
}

func TestPutUser(t *testing.T) {
	user := datastore.User{
		PublicKey: "some-key",
		DeviceId:  "some-id",
		Latitude:  5.0,
		Longitude: 6.0,
	}

	if err := user.Create(); err != nil {
		t.Fatal("Error creating user:", err)
	}

	updatedUser := datastore.User{
		Id:        user.Id,
		PublicKey: "updated-key",
	}

	jsonBytes, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatal("Error marshaling user:", err)
	}

	resp := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Id)
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(jsonBytes))

	http.DefaultServeMux.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatal("Expected 200 on user create but got", resp.Code)
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var parsedUser datastore.User
	if err := json.Unmarshal(p, &parsedUser); err != nil {
		t.Fatal(err)
	}

	if parsedUser.Id == "" {
		t.Error("Error: parsed user has no ID")
	}

	if parsedUser.PublicKey != updatedUser.PublicKey {
		t.Error("Error: expected public key", updatedUser.PublicKey, "but got", parsedUser.PublicKey)
	}
}
