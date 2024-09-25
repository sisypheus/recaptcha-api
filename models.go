package main

type Message struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email,max=100"`
	Message      string `json:"message" validate:"required,max=1000"`
	CaptchaToken string `json:"captcha_token" validate:"required"`
}
