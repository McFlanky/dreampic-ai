package util

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// checks if the password is string and meet this criteria:
// - at least 8 characters long
// - has at least one uppercase letter
// - has at least one lowercase letter
// - has at least one number
// - has at least one special character
func IsValidPassword(password string) (string, bool) {
	var (
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
		specialRunes = "!@#$%^&*"
	)

	if len(password) < 8 {
		return "Password must be at least 8 characters long", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		default:
			for _, specialChar := range specialRunes {
				if char == specialChar {
					hasSpecial = true
					break
				}
			}
		}

	}

	if !hasUpper {
		return "Password must contain at least one uppercase letter", false
	}
	if !hasLower {
		return "Password must contain at least one lowercase letter", false
	}
	if !hasNumber {
		return "Password must contain at least one number", false
	}
	if !hasSpecial {
		return "Password must contain at least one special character", false
	}

	return "Password is valid", true

}
