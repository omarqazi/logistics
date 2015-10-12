package main

import (
	"net/http"
)

const staticPath = "www"

var routes = map[string]http.Handler{
	"/locations/": LocationsController{},
	"/users/":     UserController{},
	"/":           http.FileServer(http.Dir(staticPath)),
}

func init() {
	for rule, handler := range routes {
		http.Handle(rule, http.StripPrefix(rule, handler))
	}
}
