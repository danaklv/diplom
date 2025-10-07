package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/mail"

	"gopkg.in/gomail.v2"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

func SendVerificationEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "kalykovadana3@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify your email")
	m.SetBody("text/plain", "Click the link to verify your email: http://localhost:8080/verify?code="+code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "kalykovadana3@gmail.com", "qnxq kkph idrb lvsv")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	return nil
}

func GenerateVerificationCode() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

