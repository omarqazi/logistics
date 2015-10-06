package auth

import (
	"crypto/rsa"
	"net/http"
	"time"
)

// API Tokens are valid for 15 minutes after generation
const tokenAuthorizationDuration = 15 * time.Minute

// Function request checks if the API token provided is valid
// and if not returns a 403 Forbidden returns true if request
// is authorized, false if request was blocked
func Request(w http.ResponseWriter, r *http.Request, pub *rsa.PublicKey) bool {
	token := r.Header.Get("X-API-Token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	tokenValid := TokenValid(token, tokenAuthorizationDuration, pub)
	if !tokenValid {
		http.Error(w, "Unauthorized -- Invalid API Token", 403)
		return false
	}
	return true
}
