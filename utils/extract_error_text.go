package utils

import (
    "fmt"
    "strings"
)

// ExtractErrorCodeFromText extracts the error code from a plain text error message.
// The text is expected to have the format: "Firebase Auth API error: CODE".
func ExtractErrorCodeFromText(errorMessage string) (string, error) {
    // Check if the error message starts with the expected prefix
    prefix := "Firebase Auth API error: "
    if !strings.HasPrefix(errorMessage, prefix) {
        return "", fmt.Errorf("unexpected error message format")
    }

    // Extract the error code from the message
    errorCode := strings.TrimPrefix(errorMessage, prefix)
    if errorCode == "" {
        return "", fmt.Errorf("error code not found in message")
    }

    return errorCode, nil
}