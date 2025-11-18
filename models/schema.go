package models

import (
	"encoding/json"
	"os"
	"fmt"
)

// Schema representa a estrutura completa do JSON
type Schema struct {
	TableName string  `json:"table_name"`
	Fields    []Field `json:"fields"`
}

// Field representa um campo no schema
type Field struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"` // int, string, date, text
	PrimaryKey bool       `json:"primary_key"`
	Required   bool       `json:"required"`
	Validation Validation `json:"validation"`
	Mask       string     `json:"mask"`
}

// Validation define as regras de validação
type Validation struct {
	Type       string      `json:"type"` // cpf, cnpj, email, telefone, cep, rg
	RegexRules []RegexRule `json:"regex_rules"`
}

// RegexRule define uma regra de regex customizada
type RegexRule struct {
	Pattern string `json:"pattern"`
	Message string `json:"message"`
}

// LoadSchema lê e parseia o arquivo JSON do schema
func LoadSchema(path string) (*Schema, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo schema: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal(file, &schema); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON schema: %w", err)
	}

	return &schema, nil
}
