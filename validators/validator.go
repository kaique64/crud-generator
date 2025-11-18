package validators

import (
	"net/url"
	"regexp"
	"go-crud-generator/models"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	ANY_MASK       = '*'
	CHARACTER_MASK = '#'
	DIGIT_MASK     = '9'
)

func CleanValueByMask(field models.Field, value string) string {

    if field.Mask == "" {
		return value
	}

	charsToRemove := []string{".", "-", "(", ")", "/", " ", "_"}

	cleanedValue := value

	for _, char := range charsToRemove {
		cleanedValue = strings.ReplaceAll(cleanedValue, char, "")
	}

	return cleanedValue
}

func isValidChar(inputChar rune, maskChar rune) bool {
	if maskChar == ANY_MASK {
		return true
	}

	if maskChar == DIGIT_MASK {
		return unicode.IsDigit(inputChar)
	}

	if maskChar == CHARACTER_MASK {
		return unicode.IsLetter(inputChar)
	}

	return false
}

type MaskResult struct {
	Valid     bool
	Character rune
	Consumed  bool
}

func processMaskCharacter(maskChar rune, inputChar rune) MaskResult {
	placeholders := []rune{ANY_MASK, DIGIT_MASK, CHARACTER_MASK}

	for _, placeholder := range placeholders {
		if maskChar == placeholder {
			valid := isValidChar(inputChar, maskChar)
			return MaskResult{
				Valid:     valid,
				Character: inputChar,
				Consumed:  true,
			}
		}
	}

	inputCharIsEqualToMaskChar := inputChar == maskChar

	return MaskResult{
		Valid:     true,
		Character: maskChar,
		Consumed:  inputCharIsEqualToMaskChar,
	}
}

func FormatValueByMask(mask string, value string) string {
	if mask == "" {
		return value
	}

	maskRunes := []rune(mask)
	inputRunes := []rune(value)

	result := []rune{}
	inputIndex := 0
	maskIndex := 0

	for maskIndex < len(maskRunes) && inputIndex < len(inputRunes) {
		maskResult := processMaskCharacter(maskRunes[maskIndex], inputRunes[inputIndex])

		if !maskResult.Valid {
			break
		}

		result = append(result, maskResult.Character)

		if maskResult.Consumed {
			inputIndex++
		}

		maskIndex++
	}

	return string(result)
}

func FormatDataBySchema(schema *models.Schema, data []map[string]interface{}) {
	for _, record := range data {
		for _, field := range schema.Fields {

			if field.Mask != "" {

				if rawValue, ok := record[field.Name].(string); ok {

					formattedValue := FormatValueByMask(field.Mask, rawValue)

					record[field.Name] = formattedValue
				}
			}
		}
	}
}

func FormatSingleDataBySchema(schema *models.Schema, data map[string]interface{}) {
	for _, field := range schema.Fields {
		if field.Mask != "" {
			if rawValue, ok := data[field.Name].(string); ok {
				formattedValue := FormatValueByMask(field.Mask, rawValue)
				data[field.Name] = formattedValue
			}
		}
	}
}

// ValidateData valida os dados de um formulário contra o schema e converte tipos
// Retorna um map de dados limpos e um map de erros de validação
func ValidateData(form url.Values, schema *models.Schema) (map[string]interface{}, map[string]string) {
	cleanData := make(map[string]interface{})
	errors := make(map[string]string)

	for _, field := range schema.Fields {
		value := form.Get(field.Name)

		// 1. Verificar campos obrigatórios
		if field.Required && value == "" {
			errors[field.Name] = "Campo obrigatório"
			continue
		}

		// Se não for obrigatório e estiver vazio, pulamos o resto
		if !field.Required && value == "" {
			cleanData[field.Name] = nil // Insere NULL no DB
			continue
		}

		// 2. Validações Padrão (CPF, CNPJ, etc.)
		switch field.Validation.Type {
		case "cpf":
			if !IsValidCPF(value, field.Required) {
				errors[field.Name] = "CPF inválido"
				continue
			}
			// Limpa o CPF para salvar no DB (opcional)
			// value = "somente numeros"
		case "cnpj":
			if !IsValidCNPJ(value, field.Required) {
				errors[field.Name] = "CNPJ inválido"
				continue
			}
		case "email":
			if !IsValidEmail(value, field.Required) {
				errors[field.Name] = "Email inválido"
				continue
			}
		case "cep":
			if !IsValidCEP(value, field.Required) {
				errors[field.Name] = "CEP inválido"
				continue
			}
		case "telefone":
			if !IsValidPhone(value, field.Required) {
				errors[field.Name] = "Telefone inválido"
				continue
			}
		}

		// 3. Validações de Regex Customizadas
		for _, rule := range field.Validation.RegexRules {
			matched, _ := regexp.MatchString(rule.Pattern, value)
			if !matched {
				errors[field.Name] = rule.Message
				continue // Para no primeiro erro de regex
			}
		}

		// 4. Conversão de Tipo
		// Se passou nas validações, converte para o tipo correto
		if _, hasError := errors[field.Name]; !hasError {
			switch field.Type {
			case "int":
				if field.PrimaryKey && value == "" {
					// É um create, o ID não é enviado
					continue
				}
				intVal, err := strconv.Atoi(value)
				if err != nil {
					errors[field.Name] = "Valor deve ser um número inteiro"
				} else {
					cleanData[field.Name] = intVal
				}
			case "date":
				// Tenta parsear formatos comuns (YYYY-MM-DD do HTML5 ou DD/MM/YYYY)
				dateVal, err := time.Parse("2006-01-02", value)
				if err != nil {
					dateVal, err = time.Parse("02/01/2006", value)
					if err != nil {
						errors[field.Name] = "Data inválida. Use AAAA-MM-DD"
					} else {
						cleanData[field.Name] = dateVal
					}
				} else {
					cleanData[field.Name] = dateVal
				}
			case "float":
				floatVal, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errors[field.Name] = "Valor deve ser numérico"
				} else {
					cleanData[field.Name] = floatVal
				}
			case "string", "text":
				cleanData[field.Name] = CleanValueByMask(field, value)
			default:
				cleanData[field.Name] = value
			}
		}
	}

	return cleanData, errors
}

// Helper para remover caracteres não numéricos
func justDigits(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
