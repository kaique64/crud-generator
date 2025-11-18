package validators

import "regexp"

var cepRegex = regexp.MustCompile(`^\d{8}$`)

func IsValidCEP(cep string, required bool) bool {
	if (!required) {
		return true
	}
	
	return cepRegex.MatchString(cep)
}
