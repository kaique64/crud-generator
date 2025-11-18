package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"go-crud-generator/config"
	"go-crud-generator/controllers"
	"go-crud-generator/models"
)

func main() {
	// 1. Carregar Configuração (CLI args ou ENV)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro ao carregar configuração: %v\n\n", err)
		fmt.Println("Uso:")
		fmt.Println("  ./crud-app [opções]")
		fmt.Println("\nOpções:")
		fmt.Println("  --db-host      string   Host do banco de dados (padrão: localhost)")
		fmt.Println("  --db-port      string   Porta do banco de dados (padrão: 3306)")
		fmt.Println("  --db-user      string   Usuário do banco de dados (obrigatório)")
		fmt.Println("  --db-psw       string   Senha do banco de dados")
		fmt.Println("  --db-name      string   Nome do banco de dados (obrigatório)")
		fmt.Println("  --port         string   Porta da aplicação (padrão: 8080)")
		fmt.Println("  --json-schema  string   Caminho do arquivo JSON schema (obrigatório)")
		fmt.Println("\nExemplo:")
		fmt.Println("  ./crud-app --db-host localhost --db-port 3306 --db-user root --db-psw secret --db-name mydb --port 8080 --json-schema schema.json")
		fmt.Println("\nAlternativamente, você pode usar variáveis de ambiente:")
		fmt.Println("  DB_HOST, DB_PORT, DB_USER, DB_PSW, DB_NAME, PORT, JSON_SCHEMA")
		os.Exit(1)
	}

	// Exibir configuração carregada
	log.Println("=== Configuração Carregada ===")
	log.Printf("DB Host:     %s", cfg.DBHost)
	log.Printf("DB Port:     %s", cfg.DBPort)
	log.Printf("DB User:     %s", cfg.DBUser)
	log.Printf("DB Password: %s", maskPassword(cfg.DBPassword))
	log.Printf("DB Name:     %s", cfg.DBName)
	log.Printf("App Port:    %s", cfg.Port)
	log.Printf("JSON Schema: %s", cfg.JSONSchemaPath)
	log.Println("==============================")

	// 2. Carregar Schema JSON
	schema, err := models.LoadSchema(cfg.JSONSchemaPath)
	if err != nil {
		log.Fatalf("❌ Erro ao carregar schema JSON: %v", err)
	}
	log.Println("✅ Schema JSON carregado com sucesso.")

	// 3. Conectar ao Banco de Dados
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	log.Println("✅ Conexão com MySQL estabelecida.")

	// 4. Auto-Migrate: Criar tabela se não existir
	if err := models.AutoMigrate(db, schema); err != nil {
		log.Fatalf("❌ Erro ao executar migração automática: %v", err)
	}
	log.Printf("✅ Tabela '%s' garantida.", schema.TableName)

	// 5. Inicializar Camadas
	repo := models.NewDynamicRepository(db, schema)

	// Carregar e parsear o template HTML
	tmpl, err := template.ParseFiles("views/templates/crud.html")
	if err != nil {
		log.Fatalf("❌ Erro ao parsear template: %v", err)
	}

	// 6. Configurar Controllers e Rotas
	crudController := controllers.NewCRUDController(repo, schema, tmpl)

	mux := http.NewServeMux()
	crudController.RegisterRoutes(mux)

	// Servir arquivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 7. Iniciar Servidor
	log.Printf("🚀 Servidor iniciado na porta :%s", cfg.Port)
	log.Printf("📍 Acesse: http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// maskPassword mascara a senha para exibição segura
func maskPassword(password string) string {
	if password == "" {
		return "(vazia)"
	}
	if len(password) <= 3 {
		return "***"
	}
	return password[:2] + "***" + password[len(password)-1:]
}
