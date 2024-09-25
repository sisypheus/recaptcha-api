package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

var validate = validator.New()

func SubmitHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(msg); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		sanitizer := bluemonday.UGCPolicy()
		msg.Name = sanitizer.Sanitize(msg.Name)
		msg.Email = sanitizer.Sanitize(msg.Email)
		msg.Message = sanitizer.Sanitize(msg.Message)

		remoteIP := getIPAddress(r)
		if err := VerifyRecaptcha(msg.CaptchaToken, remoteIP); err != nil {
			http.Error(w, "reCAPTCHA verification failed", http.StatusUnauthorized)
			return
		}

		if err := InsertMessage(db, msg); err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func getIPAddress(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
