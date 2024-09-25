package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func VerifyRecaptcha(token string, remoteIP string) error {
	secret := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secret == "" {
		return errors.New("reCAPTCHA secret key not set in environment variables")
	}

	verificationURL := "https://www.google.com/recaptcha/api/siteverify"

	data := url.Values{}
	data.Set("secret", secret)
	data.Set("response", token)
	data.Set("remoteip", remoteIP)

	resp, err := http.PostForm(verificationURL, data)
	if err != nil {
		return fmt.Errorf("failed to send reCAPTCHA verification request: %v", err)
	}
	defer resp.Body.Close()

	var recaptchaResp RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
		return fmt.Errorf("failed to decode reCAPTCHA response: %v", err)
	}

	if !recaptchaResp.Success {
		return fmt.Errorf("reCAPTCHA verification failed: %v", recaptchaResp.ErrorCodes)
	}

	return nil
}
