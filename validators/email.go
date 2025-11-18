package validators

import "regexp"

// (Regex simples - RFC 5322 é extremamente complexo)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail valida um email
func IsValidEmail(email string, required bool) bool {
	if (!required) {
		return true
	}

	return emailRegex.MatchString(email)
}
