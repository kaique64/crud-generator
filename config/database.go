package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
)

// InitDB inicializa e testa a conexão com o banco de dados
func InitDB(cfg *Config) (*sql.DB, error) {
	// DSN (Data Source Name)
	// Formato: "usuario:senha@tcp(host:porta)/nome_do_banco?parseTime=true"
	// parseTime=true é crucial para converter TIME e DATE do MySQL para time.Time do Go
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão SQL: %w", err)
	}

	// Testa a conexão
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("falha ao pingar banco de dados: %w", err)
	}

	// Configurações de pool de conexão
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
