package utils

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"gopkg.in/gomail.v2"
)

func ValidatePassword(password string) error {
	fmt.Println("password =", password)
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#~$%^&*()+|_]`)

	if !upper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !lower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !number.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if !special.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

func SendResetPasswordEmail(to, token string) error {
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)

	m := gomail.NewMessage()
	m.SetHeader("From", "kalykovadana3@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset Request")
	m.SetBody("text/plain",
		fmt.Sprintf(
			"We received a request to reset your password.\n\nClick the link below to set a new one (valid for 15 minutes):\n\n%s\n\nIf you didn’t request this, you can safely ignore this email.",
			resetLink,
		),
	)

	d := gomail.NewDialer("smtp.gmail.com", 587, "kalykovadana3@gmail.com", "qnxq kkph idrb lvsv")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send password reset email:", err)
		return err
	}
	return nil
}
