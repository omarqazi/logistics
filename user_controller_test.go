package main

import (
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
