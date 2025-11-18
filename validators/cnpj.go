package validators

// IsValidCNPJ valida um CNPJ
func IsValidCNPJ(cnpj string, required bool) bool {
	if (!required) {
		return true
	}
	
	cnpj = justDigits(cnpj)
	if len(cnpj) != 14 {
		return false
	}
	if allSameDigits(cnpj) {
		return false
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	d1 := calculateCNPJDigit(cnpj[:12], weights1)
	d2 := calculateCNPJDigit(cnpj[:13], weights2)

	return cnpj[12] == d1 && cnpj[13] == d2
}

func calculateCNPJDigit(doc string, weights []int) uint8 {
	sum := 0
	for i, r := range doc {
		sum += (int(r) - '0') * weights[i]
	}
	rem := sum % 11
	if rem < 2 {
		return '0'
	}
	return uint8(11-rem) + '0'
}
