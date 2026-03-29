package utils

import (
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func Validate(r interface{}) validator.ValidationErrors {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func PasswordValidation(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 12 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func EmailValidation(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func ParseDate(s *string) (*time.Time, error) {
    if s == nil || *s == "" {
        return nil, nil
    }

    formats := []string{
        "2006-01-02",
        time.RFC3339,
        "2006-01-02T15:04:05Z",
    }

    for _, format := range formats {
        t, err := time.Parse(format, *s)
        if err == nil {
            return &t, nil
        }
    }

    return nil, fmt.Errorf("cannot parse date: %s", *s)
}
