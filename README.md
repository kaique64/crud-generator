# Gerador de CRUD Dinâmico em Go

Esta é uma aplicação web full-stack em Go que gera automaticamente uma interface web CRUD (Create, Read, Update, Delete) completa com base em um schema JSON.

A aplicação utiliza o padrão MVC, MySQL como banco de dados e TailwindCSS para o frontend.

## 🚀 Funcionalidades

* **Geração Dinâmica:** A aplicação lê um `schema.json` na inicialização.
* **Auto-Migração:** Cria automaticamente a tabela no MySQL (usando `CREATE TABLE IF NOT EXISTS`) com base no schema.
* **CRUD Completo:** Interface web para Criar, Listar (com paginação e busca), Atualizar e Excluir registros.
* **Validação Backend:** Validação robusta no lado do servidor (Obrigatório, CPF, CNPJ, Email, Regex) antes de salvar no banco.
* **Validação Frontend:** Validação e máscaras de entrada (CPF, Telefone, CEP) no lado do cliente.
* **Arquitetura Limpa:** Padrão MVC com separação clara de responsabilidades.
* **Segurança:** Utiliza *prepared statements* para prevenir SQL Injection e `html/template` para prevenir XSS.

## 🛠️ Stack

* **Backend:** Go (stdlib `net/http`)
* **Banco de Dados:** MySQL (5.7+ e 8.0+)
* **Driver DB:** `go-sql-driver/mysql`
* **Frontend:** HTML, TailwindCSS (via CDN), Vanilla JavaScript

## ⚙️ Configuração e Execução

1.  **Pré-requisitos:**
    * Go (v1.20+)
    * MySQL (5.7 ou 8.0)

2.  **Setup do Banco:**
    * Crie um banco de dados no seu MySQL. Ex: `CREATE DATABASE meu_crud_db;`

3.  **Schema:**
    * Crie seu arquivo `schema.json` (veja o exemplo na especificação) e salve-o (ex: `./schema.json`).

4.  **Dependências:**
    * Execute `go mod tidy` para baixar a dependência do driver MySQL.

5.  **Variáveis de Ambiente:**
    * A aplicação é configurada via variáveis de ambiente. Você pode exportá-las ou usar um arquivo `.env` (com `source .env`).

    ```bash
    export DB_HOST="localhost"
    export DB_PORT="3306"
    export DB_NAME="meu_crud_db"
    export DB_USER="seu_usuario_mysql"
    export DB_PSW="sua_senha_mysql"
    export JSON_SCHEMA="./schema.json" # Caminho para seu schema
    export PORT="8080"
    ```

6.  **Compilar e Executar:**

    ```bash
    # Compilar
    go build -o crud-app .

    # Executar (com as variáveis de ambiente carregadas)
    ./crud-app
    ```

7.  **Acessar:**
    * Abra seu navegador e acesse `http://localhost:8080`.

## 🏛️ Arquitetura

* `main.go`: Ponto de entrada, "cola" da aplicação.
* `config/`: Carregamento de env vars (`config.go`) e conexão com DB (`database.go`).
* `models/`:
    * `schema.go`: Structs e parser do JSON.
    * `migration.go`: Lógica do `CREATE TABLE`.
    * `repository.go`: O "Model" dinâmico. Constrói queries SQL seguras.
* `controllers/`:
    * `crud_controller.go`: Os "Controllers" (handlers HTTP). Gerencia o request, chama o repositório/validador e renderiza a view.
* `validators/`: Pacote com toda a lógica de validação de dados (CPF, CNPJ, Email, etc.).
* `views/templates/`:
    * `crud.html`: O "View". Template HTML único que se renderiza dinamicamente.
* `static/js/`:
    * `main.js`: JavaScript do frontend para máscaras, validação e modo de edição.

## ⚠️ Limitações e Próximos Passos

Como um sistema de *scaffolding* em tempo real, esta prova de conceito é robusta, mas pode ser estendida:

* **Segurança (CSRF):** Implementar tokens Anti-CSRF para proteger contra ataques de falsificação de solicitação.
* **Tipos de Campo:** Suportar mais tipos de campo (ex: `<select>`, `<textarea>`, `checkbox`).
* **Soft Delete:** Adicionar a lógica de "soft delete" (baseado em uma flag no schema).
* **Relações:** Suportar chaves estrangeiras (ex: `belongs_to`), o que aumentaria drasticamente a complexidade.
