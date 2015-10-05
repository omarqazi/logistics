package main

import (
	"net/http"
	"testing"
)

// TestRouter tests to make sure all the routes
// are mapped properly.
func TestRouter(t *testing.T) {
	for pattern, handler := range routes {
		req, err := http.NewRequest("GET", pattern, nil)
		if err != nil {
			t.Error(err)
		}

		rh, pat := http.DefaultServeMux.Handler(req)
		if rh != handler {
			t.Error("Error: expected handler", handler, "but got", rh, "for route", pattern, "(matched ", pat, ")")
		}
	}
}
