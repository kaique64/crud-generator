package validators

// IsValidCPF valida um CPF
func IsValidCPF(cpf string, required bool) bool {
	if (!required) {
		return true
	}
	
	cpf = justDigits(cpf)
	if len(cpf) != 11 {
		return false
	}
	if allSameDigits(cpf) {
		return false
	}

	d1 := calculateDigit(cpf[:9], 10)
	d2 := calculateDigit(cpf[:10], 11)

	return cpf[9] == d1 && cpf[10] == d2
}

func allSameDigits(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func calculateDigit(doc string, weight int) uint8 {
	sum := 0
	for _, r := range doc {
		sum += (int(r) - '0') * weight
		weight--
	}
	rem := sum % 11
	if rem < 2 {
		return '0'
	}
	return uint8(11-rem) + '0'
}
