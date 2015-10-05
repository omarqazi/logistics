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
	} else if r.Method == "POST" {
		postUser(w, r)
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

	user.Latitude = 0.0
	user.Longitude = 0.0
	user.UpdatedAt = user.CreatedAt
	user.DeviceId = "[redacted]"

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user:", err)
		http.Error(w, "Error serializing data", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintln(w, string(bytes))
}

func postUser(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user datastore.User
	if err := dec.Decode(&user); err != nil {
		log.Println("Error decoding user:", err)
		http.Error(w, "Error decoding user", 400)
		return
	}

	if err := user.Create(); err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Error creating user", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(user); err != nil {
		log.Println("Error marshaling saved user:", err)
		http.Error(w, "Error marshaling user data", 500)
		return
	}
	return
}
