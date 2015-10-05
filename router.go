package main

import (
	"net/http"
)

var routes = map[string]http.Handler{
	"/locations/": LocationsController{},
	"/users/":     UserController{},
}

func init() {
	for rule, handler := range routes {
		http.Handle(rule, http.StripPrefix(rule, handler))
	}
}
