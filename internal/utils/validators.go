package utils

import (
	"regexp"
	"strings"
	"time"
)

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	return re.MatchString(strings.ToLower(email))
}

func IsValidUsername(username string) bool {
	return len(username) >= 3 && len(username) <= 30
}

func IsValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 20
}

// Valida si la persona tiene entre 13 y 120 años (ya más es no creo)
func IsValidDate(date string) bool {
	if date == "" {
		return false
	}

	birthDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	today := time.Now()

	age := today.Year() - birthDate.Year()
	hasHadBirthdayThisYear := today.Month() > birthDate.Month() ||
		(today.Month() == birthDate.Month() && today.Day() >= birthDate.Day())

	exactAge := age
	if !hasHadBirthdayThisYear {
		exactAge = age - 1
	}

	return exactAge >= 13 && exactAge <= 120
}
