package utils

import (
	"errors"
	"regexp"
)

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	// Должен начинаться с буквы
	startWithLetter := regexp.MustCompile(`^[A-Za-z]`)
	if !startWithLetter.MatchString(username) {
		return errors.New("username must start with a letter")
	}

	// Разрешённые символы: буквы, цифры, _
	allowed := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	if !allowed.MatchString(username) {
		return errors.New("username can only contain letters, digits, and underscore")
	}

	// Не только цифры
	onlyDigits := regexp.MustCompile(`^[0-9]+$`)
	if onlyDigits.MatchString(username) {
		return errors.New("username cannot contain only digits")
	}

	return nil
}
