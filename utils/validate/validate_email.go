package utils_validate

import (
	"net"
	"strings"
	"unicode"
)

// IsEmailValid validates the email address using string operations
func IsEmailValid(email string) bool {
	// Check the overall length of the email
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Split the email into local part and domain part
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:]

	// Validate the local part
	if !validateLocalPart(localPart) {
		return false
	}

	// Validate the domain part
	if !validateDomainPart(domainPart) {
		return false
	}

	return true
}

func validateLocalPart(localPart string) bool {
	if len(localPart) < 1 || len(localPart) > 64 {
		return false
	}

	// The local part can contain letters, digits, and special characters
	allowedSpecialChars := "!#$%&'*+/=?^_`{|}~."
	lastChar := ' '
	inQuotes := false

	for i, char := range localPart {
		if !(unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune(allowedSpecialChars, char)) {
			if char == '\\' && inQuotes {
				// Allow escaped characters within quotes
				if i+1 < len(localPart) {
					i++
					continue
				}
				return false
			}
			if char == '"' {
				if i == 0 || (i > 0 && localPart[i-1] == '\\') {
					inQuotes = !inQuotes
					continue
				}
				return false
			}
			if !inQuotes {
				return false
			}
		}

		// Dot cannot appear consecutively
		if char == '.' && lastChar == '.' {
			return false
		}

		lastChar = char
	}

	// Check if local part consists only of special characters
	if strings.Trim(localPart, allowedSpecialChars) == "" {
		return false
	}

	return true
}

func validateDomainPart(domainPart string) bool {
	if len(domainPart) < 1 || len(domainPart) > 253 {
		return false
	}

	// Handle IP address
	if strings.HasPrefix(domainPart, "[") && strings.HasSuffix(domainPart, "]") {
		ip := domainPart[1 : len(domainPart)-1]
		if strings.HasPrefix(ip, "IPv6:") {
			ip = ip[5:]
		}
		if net.ParseIP(ip) == nil {
			return false
		}
		return true
	}

	// Split the domain into labels
	labels := strings.Split(domainPart, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if len(label) < 1 || len(label) > 63 {
			return false
		}

		// Labels must start and end with a letter or digit
		if !(unicode.IsLetter(rune(label[0])) || unicode.IsDigit(rune(label[0]))) {
			return false
		}
		if !(unicode.IsLetter(rune(label[len(label)-1])) || unicode.IsDigit(rune(label[len(label)-1]))) {
			return false
		}

		// Labels can contain letters, digits, and hyphens
		for _, char := range label {
			if !(unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-') {
				return false
			}
		}
	}

	// The TLD (last label) must be at least 2 characters
	tld := labels[len(labels)-1]
	if len(tld) < 2 {
		return false
	}

	return true
}
