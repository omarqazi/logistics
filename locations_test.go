package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLocationsController(t *testing.T) {
	resp := httptest.NewRecorder()
	uri := "/locations/whatever"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	locationsController := LocationsController{}
	locationsController.ServeHTTP(resp, req)

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else {
		if !strings.Contains(string(p), "locations controller") {
			t.Error("Locations controller response does not contain the phrase locations controller")
		}
	}
}
