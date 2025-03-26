package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"senatus/components"
	"senatus/models"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Global repository
var todoRepo models.Repository

func main() {
	// Setup SQLite database in the data directory
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dataDir, "senatus.db")
	repo, err := models.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize repository
	todoRepo = repo

	// Ensure the database has initial data
	if sqliteRepo, ok := todoRepo.(*models.SQLiteRepository); ok {
		if err := sqliteRepo.SeedInitialData(); err != nil {
			log.Printf("Warning: Failed to seed initial data: %v", err)
		}

		// Ensure we close the database connection when the server shuts down
		defer sqliteRepo.Close()
	}

	// Create new Go 1.22 ServeMux
	mux := http.NewServeMux()

	// Static file handler
	staticHandler := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	// HTML routes
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /view/{id}", viewTodoHandler)
	mux.HandleFunc("GET /new", newTodoHandler)
	mux.HandleFunc("GET /edit/{id}", editTodoHandler)
	mux.HandleFunc("GET /filter/status/{status}", filterByStatusHandler)
	mux.HandleFunc("GET /filter/priority/{priority}", filterByPriorityHandler)

	// API routes
	mux.HandleFunc("POST /todo", createTodoHandler)
	mux.HandleFunc("GET /todo/{id}", getTodoHandler)
	mux.HandleFunc("PUT /todo/{id}", updateTodoHandler)
	mux.HandleFunc("DELETE /todo/{id}", deleteTodoHandler)
	mux.HandleFunc("PUT /todo/{id}/status", updateTodoStatusHandler)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// homeHandler renders the home page with all todos
func homeHandler(w http.ResponseWriter, r *http.Request) {
	todos := todoRepo.GetAll()
	render(w, r, components.HomePage(todos))
}

// viewTodoHandler renders a single todo detail view
func viewTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	todo, found := todoRepo.GetByID(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	render(w, r, components.TodoDetail(todo))
}

// newTodoHandler renders the new todo form
func newTodoHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, components.NewTaskForm(models.Todo{
		Priority: 3, // Default to medium priority
	}))
}

// editTodoHandler renders the edit todo form
func editTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	todo, found := todoRepo.GetByID(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	render(w, r, components.EditTaskForm(todo))
}

// filterByStatusHandler shows todos filtered by status
func filterByStatusHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	status := models.Status(parts[3])
	todos := todoRepo.GetByStatus(status)

	render(w, r, components.HomePage(todos))
}

// filterByPriorityHandler shows todos filtered by priority
func filterByPriorityHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	priority, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}

	todos := todoRepo.FilterByPriority(priority)

	render(w, r, components.HomePage(todos))
}

// createTodoHandler creates a new todo
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	priority, _ := strconv.Atoi(r.FormValue("priority"))

	todo := models.Todo{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Priority:    priority,
		Sponsor:     r.FormValue("sponsor"),
		Status:      models.Proposed,
	}

	// Save the new todo
	todoRepo.Create(todo)

	// Redirect to the home page after successful creation
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

// getTodoHandler returns a single todo as JSON
func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	todo, found := todoRepo.GetByID(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// updateTodoHandler updates an existing todo
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the existing todo
	existingTodo, found := todoRepo.GetByID(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	priority, _ := strconv.Atoi(r.FormValue("priority"))

	// Update the todo with form values
	existingTodo.Title = r.FormValue("title")
	existingTodo.Description = r.FormValue("description")
	existingTodo.Priority = priority
	existingTodo.Sponsor = r.FormValue("sponsor")

	statusValue := r.FormValue("status")
	if statusValue != "" {
		existingTodo.Status = models.Status(statusValue)
	}

	// Save the updated todo
	_, updated := todoRepo.Update(existingTodo)
	if !updated {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	// Redirect to the home page or detail page after successful update
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

// deleteTodoHandler deletes a todo
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)

	deleted := todoRepo.Delete(id)
	if !deleted {
		http.Error(w, "Todo not found or could not be deleted", http.StatusNotFound)
		return
	}

	// If the request is from the detail page and requests a redirect
	if r.Header.Get("HX-Push-Url") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	// For HTMX delete requests from the list view, just return 200 OK
	w.WriteHeader(http.StatusOK)
}

// updateTodoStatusHandler updates just the status of a todo
func updateTodoStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(strings.TrimSuffix(r.URL.Path, "/status"))

	// Get the existing todo
	existingTodo, found := todoRepo.GetByID(id)
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Parse the JSON request body
	var requestBody struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the status
	existingTodo.Status = models.Status(requestBody.Status)

	// Save the updated todo
	updatedTodo, updated := todoRepo.Update(existingTodo)
	if !updated {
		http.Error(w, "Failed to update todo status", http.StatusInternalServerError)
		return
	}

	// Return the updated status change buttons
	render(w, r, components.StatusChangeButtons(updatedTodo))
}

// render is a helper function to render a templ component
func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
