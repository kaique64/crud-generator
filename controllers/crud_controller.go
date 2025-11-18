package controllers

import (
	"html/template"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"go-crud-generator/models"
	"go-crud-generator/validators"
	"strconv"
	"time"
)

const defaultPageLimit = 10

// CRUDController gerencia as rotas e handlers do CRUD
type CRUDController struct {
	repo   *models.DynamicRepository
	schema *models.Schema
	tmpl   *template.Template
}

// NewCRUDController cria uma nova instância do controller
func NewCRUDController(repo *models.DynamicRepository, schema *models.Schema, tmpl *template.Template) *CRUDController {
	return &CRUDController{
		repo:   repo,
		schema: schema,
		tmpl:   tmpl,
	}
}

// RegisterRoutes registra as rotas no mux
func (c *CRUDController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.handleList)
	mux.HandleFunc("/create", c.handleCreate)
	mux.HandleFunc("/update", c.handleUpdate) // Usará /update?id=...
	mux.HandleFunc("/delete", c.handleDelete) // Usará /delete?id=...
	mux.HandleFunc("/get", c.handleGetByID)   // Rota AJAX para editar
}

// TemplateData é a estrutura de dados passada para o template HTML
type TemplateData struct {
	Schema       *models.Schema
	Data         []map[string]interface{}
	Errors       map[string]string
	FormData     map[string]string // Para repopular o form em caso de erro
	SearchTerm   string
	Pagination   Pagination
	CurrentTime  int64 // Para cache-busting de estáticos
	SuccessMessage string
    SchemaColspan int // <- ADICIONE ESTA LINHA
}

// Pagination contém dados para a paginação
type Pagination struct {
	CurrentPage int
	TotalPages  int
	TotalRecords int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
}

// handleList exibe a página principal com a lista e o formulário
func (c *CRUDController) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	pageStr := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	data, totalRecords, err := c.repo.FindAll(page, defaultPageLimit, search)
	if err != nil {
		log.Printf("Erro ao buscar dados: %v", err)
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(defaultPageLimit)))
	pagination := Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalRecords: totalRecords,
		HasPrev:     page > 1,
		PrevPage:    page - 1,
		HasNext:     page < totalPages,
		NextPage:    page + 1,
	}

	validators.FormatDataBySchema(c.schema, data)
	
	templateData := TemplateData{
		Schema:     c.schema,
		Data:       data,
		SearchTerm: search,
		Pagination: pagination,
		CurrentTime: time.Now().Unix(),
		SchemaColspan: len(c.schema.Fields) + 1,
	}
	
	c.renderTemplate(w, templateData)
}

// handleCreate processa a submissão do formulário de criação
func (c *CRUDController) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}
	
	// Validar e converter dados
	data, validationErrors := validators.ValidateData(r.PostForm, c.schema)

	if len(validationErrors) > 0 {
		// Recarregar a página com erros
		c.reloadPageWithErrors(w, r, validationErrors, r.PostForm)
		return
	}

	// Inserir no banco
	_, err := c.repo.Create(data)
	if err != nil {
		log.Printf("Erro ao criar registro: %v", err)
		validationErrors["_form"] = "Erro interno ao salvar. Verifique se os dados estão corretos."
		c.reloadPageWithErrors(w, r, validationErrors, r.PostForm)
		return
	}

	// Redireciona para a home (pode adicionar ?success=true)
	http.Redirect(w, r, "/", http.StatusFound)
}

// handleUpdate processa a submissão do formulário de edição
func (c *CRUDController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID ausente", http.StatusBadRequest)
		return
	}

	data, validationErrors := validators.ValidateData(r.PostForm, c.schema)
	if len(validationErrors) > 0 {
		c.reloadPageWithErrors(w, r, validationErrors, r.PostForm)
		return
	}

	var pkValue interface{}
	for _, field := range c.schema.Fields {
		if field.PrimaryKey {
			pkValue = data[field.Name]
			break
		}
	}

	if err := c.repo.Update(pkValue, data); err != nil {
		log.Printf("Erro ao atualizar registro: %v", err)
		validationErrors["_form"] = "Erro interno ao atualizar."
		c.reloadPageWithErrors(w, r, validationErrors, r.PostForm)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// handleDelete processa a exclusão de um item (via POST para segurança)
func (c *CRUDController) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID ausente", http.StatusBadRequest)
		return
	}

	// Converte ID para o tipo da PK (assumindo int por simplicidade)
	// Em um sistema real, leríamos o tipo da PK do schema
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := c.repo.Delete(idInt); err != nil {
		log.Printf("Erro ao deletar registro: %v", err)
		http.Error(w, "Erro ao deletar registro", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// handleGetByID é usado pelo AJAX para popular o formulário de edição
func (c *CRUDController) handleGetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID ausente", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	data, err := c.repo.FindByID(idInt)
	
	validators.FormatSingleDataBySchema(c.schema, data)

	if err != nil {
		log.Printf("Erro ao buscar por ID: %v", err)
		http.Error(w, "Registro não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


// renderTemplate renderiza o template HTML com os dados fornecidos
func (c *CRUDController) renderTemplate(w http.ResponseWriter, data TemplateData) {
	err := c.tmpl.ExecuteTemplate(w, "crud.html", data)
	if err != nil {
		log.Printf("Erro ao renderizar template: %v", err)
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
	}
}

// reloadPageWithErrors recarrega a página de lista, injetando os erros de validação
func (c *CRUDController) reloadPageWithErrors(w http.ResponseWriter, r *http.Request, errors map[string]string, formData map[string][]string) {
	pageStr := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	data, totalRecords, err := c.repo.FindAll(page, defaultPageLimit, search)
	if err != nil {
		log.Printf("Erro ao buscar dados: %v", err)
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(defaultPageLimit)))
	pagination := Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		PrevPage:    page - 1,
		HasNext:     page < totalPages,
		NextPage:    page + 1,
	}

	// Converte url.Values (map[string][]string) para map[string]string
	simpleFormData := make(map[string]string)
	for k, v := range formData {
		if len(v) > 0 {
			simpleFormData[k] = v[0]
		}
	}

	templateData := TemplateData{
		Schema:     c.schema,
		Data:       data,
		SearchTerm: search,
		Pagination: pagination,
		Errors:     errors,
		FormData:   simpleFormData,
		CurrentTime: time.Now().Unix(),
		SchemaColspan: len(c.schema.Fields) + 1,
	}

	w.WriteHeader(http.StatusBadRequest) // Indica que foi um request inválido
	c.renderTemplate(w, templateData)
}
