package main

import (
	"fmt"
	"net/http"
)

type LocationsController struct {
}

func (l LocationsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the locations controller")
}
