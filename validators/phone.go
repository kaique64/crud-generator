package validators

import "regexp"

var phoneRegex = regexp.MustCompile(`^\d{10,11}$`)

func IsValidPhone(phone string, required bool) bool {
	if (!required) {
		return true
	}

	return phoneRegex.MatchString(phone)
}
