package handlers

import (
	"cryonics/internal/auth"
	"fmt"
	"log"
	"net/http"
)

type Acc struct {
	Token string
	UID   string
}

var UserChan = make(chan Acc, 1)

func ReceiveToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.FormValue("idToken")
	if token == "" {
		http.Error(w, "Missing idToken", http.StatusBadRequest)
		return
	}

	log.Printf("Received ID token from browser: %s...", token[:20])

	claims, err := auth.VerifyFirebaseIDToken(r.Context(), token)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	log.Println("Token verified successfully!")
	log.Printf("User ID: %s", claims.Subject)

	if identities, ok := claims.Firebase["identities"].(map[string]interface{}); ok {
		if emails, ok := identities["email"].([]interface{}); ok && len(emails) > 0 {
			log.Printf("Email: %s", emails[0])
		}
	}

	UserChan <- Acc{Token: token, UID: claims.Subject}

	fmt.Fprintf(w, "Token verified! Welcome user: %s", claims.Subject)
}
