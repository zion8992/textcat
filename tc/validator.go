package tc

import (
	"regexp"
	"fmt"
	"unicode/utf8"
)

var validChars = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func IsValidUsername(s string) bool {
	return validChars.MatchString(s)
}

func ValidateMessage(msg string, maxChars int) (string, error) {
	// Sanitize by removing < and >
	re := regexp.MustCompile(`[<>]`)
	safe := re.ReplaceAllString(msg, "")

	// Count runes (not bytes)
	runeCount := utf8.RuneCountInString(safe)
	if runeCount > maxChars {
		return "", fmt.Errorf("message too long: %d runes (max %d)", runeCount, maxChars)
	}

	return safe, nil
}
