package main

import (
	"encoding/json"
	"fmt"
	"github.com/omarqazi/logistics/auth"
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
	} else if r.Method == "PUT" {
		putUser(w, r)
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
	user.LocatedAt = user.CreatedAt
	user.UpdatedAt = user.CreatedAt

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

func putUser(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var user datastore.User
	if err := dec.Decode(&user); err != nil {
		log.Println("Error decoding user:", err)
		http.Error(w, "Error decoding user", 400)
		return
	}

	dbUser, err := datastore.GetUser(user.Id)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	rsaKey, _ := dbUser.RSAKey()
	if ok := auth.Request(w, r, rsaKey); !ok {
		return
	}

	if user.PublicKey == "" {
		user.PublicKey = dbUser.PublicKey
	}

	if user.DeviceId == "" {
		user.DeviceId = dbUser.DeviceId
	}

	if user.Latitude == 0.0 && user.Longitude == 0.0 {
		user.Latitude = dbUser.Latitude
		user.Longitude = dbUser.Longitude
	}

	if user.LatestLocationId == "" {
		user.LatestLocationId = dbUser.LatestLocationId
	}

	if user.LocatedAt.Before(dbUser.LocatedAt) {
		user.LocatedAt = dbUser.LocatedAt
	}

	if err := user.Update(); err != nil {
		log.Println("Error updating user", user, err)
		http.Error(w, "Error updating user", 500)
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(user); err != nil {
		log.Println("Error marshaling saved user:", err)
		http.Error(w, "Error marshaling user data", 500)
	}
	return
}
