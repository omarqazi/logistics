package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
}
