package main

import (
	"net/http"
)

var routes = map[string]http.Handler{
	"/locations/": LocationsController{},
}

func init() {
	for rule, handler := range routes {
		http.Handle(rule, handler)
	}
}
