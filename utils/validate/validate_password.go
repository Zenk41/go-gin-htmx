package utils_validate

import (
	"strings"
	"unicode"
)

// IsPasswordValid validates the password based on given criteria
func IsPasswordValid(password string) (bool, string) {
	const minLen = 8
	const maxLen = 64

	if len(password) < minLen || len(password) > maxLen {
		return false, "Need Min 8 & Max 64 Character length"
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?/"

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "Need At least One Upper Character"
	}
	if !hasLower {
		return false, "Need At least One Lower Character"
	}
	if !hasDigit {
		return false, "Need At least One Numeric"
	}
	if !hasSpecial {
		return false, "Need At least One Special Character"
	}

	return true, ""
}
