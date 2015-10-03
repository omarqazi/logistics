package main

import (
	"fmt"
	"logistics/datastore"
	"net/http"
)

type LocationsController struct {
}

func (l LocationsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the locations controller", datastore.Redis)
}
