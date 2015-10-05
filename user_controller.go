package main

import (
	"encoding/json"
	"fmt"
	"github.com/omarqazi/logistics/datastore"
	"log"
	"net/http"
)

type UserController struct {
}

func (uc UserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getUser(w, r)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	user, err := datastore.GetUser(r.URL.Path)
	if err != nil {
		log.Println("Error getting user:", err)
		http.Error(w, "Error getting user", 500)
		return
	}
	if user == nil {
		http.Error(w, "User does not exist", 404)
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user:", err)
		http.Error(w, "Error serializing data", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintln(w, string(bytes))
}
