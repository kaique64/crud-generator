package models

import (
	"database/sql"
	"fmt"
	"strings"
)

// AutoMigrate cria a tabela no banco de dados com base no schema, se ela não existir
func AutoMigrate(db *sql.DB, schema *Schema) error {
	query := buildCreateTableQuery(schema)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("falha ao executar CREATE TABLE: %w. Query: %s", err, query)
	}
	return nil
}

// buildCreateTableQuery constrói a string da query SQL
func buildCreateTableQuery(schema *Schema) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", schema.TableName))

	var primaryKey string
	definitions := []string{}

	for _, field := range schema.Fields {
		sqlType := mapJSONTypeToSQL(field.Type)
		definition := fmt.Sprintf("  %s %s", field.Name, sqlType)

		if field.PrimaryKey {
			if field.Type == "int" {
				definition += " AUTO_INCREMENT"
			}
			primaryKey = field.Name
		} else if field.Required {
			definition += " NOT NULL"
		} else {
			definition += " NULL"
		}

		definitions = append(definitions, definition)
	}

	sb.WriteString(strings.Join(definitions, ",\n"))

	if primaryKey != "" {
		sb.WriteString(fmt.Sprintf(",\n  PRIMARY KEY (%s)", primaryKey))
	}

	sb.WriteString("\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")

	return sb.String()
}

// mapJSONTypeToSQL traduz tipos do JSON para tipos SQL compatíveis com MySQL 5.7+ e 8.0+
func mapJSONTypeToSQL(jsonType string) string {
	switch jsonType {
	case "int":
		return "INT"
	case "string":
		return "VARCHAR(255)"
	case "text": // Para campos maiores
		return "TEXT"
	case "date":
		return "DATE"
	case "datetime":
		return "DATETIME"
	case "float":
		return "DECIMAL(10, 2)" // Padrão genérico
	default:
		return "VARCHAR(255)"
	}
}
