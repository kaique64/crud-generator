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
    
    # Ou pode rodar dessa maneira, caso esteja no Linux
    DB_HOST=localhost DB_PORT=3306 DB_NAME=crud_app DB_USER=root DB_PSW=root JSON_SCHEMA=schema.json PORT=8081 ./crud-app
    ```

7.  **Acessar:**
    * Abra seu navegador e acesse `http://localhost:8080`.

8. **Executar com WINDOWS**
```bash
GOOS=windows GOARCH=amd64 go build -o crud-app.exe main.go

./crud-app.exe --db-host localhost --db-port 3306 --db-user root --db-psw root --db-name crud_app --port 8081 --json-schema schema.json
```

# 📖 Guia de Configuração: `schema.json`

Este arquivo `schema.json` é o coração do sistema, definindo a estrutura da tabela no banco de dados e as regras de exibição e validação no frontend.

## Estrutura Básica

O schema é composto por um objeto principal que contém o nome da tabela (`TableName`) e uma lista de campos (`Fields`).

```json
{
    "table_name": "nome_da_tabela",
    "fields": [
        {
            // Definição do campo 1
        },
        {
            // Definição do campo 2
        }
    ]
}
```
-----

## Detalhe dos Campos (`Fields`)

Cada objeto dentro da lista `Fields` define uma coluna no banco de dados e suas propriedades na aplicação:

| Propriedade | Tipo | Obrigatório | Descrição | Exemplo de Valor |
| :--- | :--- | :--- | :--- | :--- |
| `name` | string | Sim | Nome da coluna no banco de dados. Deve ser único. | `"cpf"`, `"nome"`, `"id"` |
| `type` | string | Sim | Tipo de dado (usado para renderização do input e tipagem no Go/Gorm). | `"string"`, `"int"`, `"date"`, `"text"` |
| `primary_key` | bool | Não | Define se o campo é a chave primária da tabela. | `true` |
| `required` | bool | Não | Define se o campo é obrigatório (validação de frontend e backend). | `true` |
| `mask` | string | Não | Máscara de formatação para o frontend (IMask.js). **Ver Regras de Máscara abaixo.** | `"999.999.999-99"` |
| `validation` | objeto | Não | Objeto que define o tipo de validação de frontend e backend. | Ver **Regras de Validação** |

-----

## 🔑 Regras de Máscara (`Mask`)

Use esta propriedade para formatar a entrada de dados no formulário (frontend). A validação (backend) receberá apenas o valor puro.

| Símbolo | Significado | Exemplo de Uso | Resultado Esperado |
| :--- | :--- | :--- | :--- |
| **`9`** | **Dígito** (0-9). | `"99999-999"` | `12345-678` (CEP) |
| **`#`** | **Caractere** (Letra A-Z, a-z). | `"####-999"` | `ABCD-123` |
| **`*`** | **Qualquer tipo** (Dígito, Letra, Símbolo). | `"AA*-99"` | `AAx-12` |
| **Outros** | Caracteres fixos (pontuação). | N/A | Caracteres fixos (ex: `.`, `-`, `/`, `(`). |

### Exemplos de Máscara:

| Campo | Máscara |
| :--- | :--- |
| CPF | `"999.999.999-99"` |
| CNPJ | `"99.999.999/9999-99"` |
| Placa | `"###-9999"` |
| Telefone (Dinâmico) | `"(99) 99999-9999"` (O JS lida com 10/11 dígitos automaticamente) |

-----

## 🔎 Regras de Validação (`Validation`)

Use este objeto para aplicar validações específicas no campo.

```json
"validation": {
    "type": "cpf" // O nome da validação (usado no switch/case do Go e JS)
    // "regex": "^[a-zA-Z]+$", // (Opcional, para validação por expressão regular)
}
```

| `calidation.type` | Descrição |
| :--- | :--- |
| `"cpf"` | Validação de CPF (dígito verificador). |
| `"cnpj"` | Validação de CNPJ (dígito verificador). |
| `"email"` | Validação de formato de email (`@`, `.com`, etc.). |
| `"cep"` | Validação de CEP (8 dígitos). |
| `"telefone"` | Validação de telefone (10 ou 11 dígitos). |

-----

## Exemplo Completo de `schema.json`

```json
{
    "table_name": "clientes",
    "fields": [
        {
            "name": "id",
            "type": "int",
            "primary_key": true
        },
        {
            "name": "nome",
            "type": "string",
            "required": true
        },
        {
            "name": "cpf",
            "type": "string",
            "required": true,
            "Mask": "999.999.999-99",
            "validation": {
                "type": "cpf"
            }
        },
        {
            "name": "email",
            "type": "string",
            "required": false,
            "validation": {
                "type": "email"
            }
        },
        {
            "name": "cep",
            "type": "string",
            "required": false,
            "mask": "99999-999",
            "validation": {
                "type": "cep"
            }
        }
    ]
}
```

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
