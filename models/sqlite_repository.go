package models

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository implements Repository with a SQLite database
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Ensure the directory exists
	if dbPath != ":memory:" {
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, err
		}
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Set pragmas for better performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;"); err != nil {
		db.Close()
		return nil, err
	}

	// Create the todos table if it doesn't exist
	if err := initializeDB(db); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Initialize the database with required tables
func initializeDB(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		priority INTEGER NOT NULL,
		sponsor TEXT NOT NULL
	);
	`

	_, err := db.Exec(createTableSQL)
	return err
}

// SeedInitialData inserts sample data if the table is empty
func (r *SQLiteRepository) SeedInitialData() error {
	// Check if the table is empty
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	if err != nil {
		return err
	}

	// If there are already records, don't seed
	if count > 0 {
		return nil
	}

	// Sample data for seeding
	initialTodos := []Todo{
		{
			ID:          uuid.NewString(),
			Title:       "Plan Aqueduct Expansion",
			Description: "Propose plans for extending the aqueduct system to the eastern provinces",
			Status:      Proposed,
			CreatedAt:   time.Now().Add(-72 * time.Hour),
			UpdatedAt:   time.Now().Add(-72 * time.Hour),
			Priority:    2,
			Sponsor:     "Senator Marcus Aquillius",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Review Military Budget",
			Description: "Examine the funding allocation for legions stationed at the northern frontier",
			Status:      Approved,
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
			Priority:    1,
			Sponsor:     "Consul Lucius Cornelius",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Organize Grain Distribution",
			Description: "Coordinate the monthly grain distribution to citizens",
			Status:      InProgress,
			CreatedAt:   time.Now().Add(-96 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
			Priority:    3,
			Sponsor:     "Aedile Publius Clodius",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Commission Temple Repairs",
			Description: "Allocate funds for restoration of the Temple of Jupiter",
			Status:      Completed,
			CreatedAt:   time.Now().Add(-120 * time.Hour),
			UpdatedAt:   time.Now().Add(-6 * time.Hour),
			Priority:    4,
			Sponsor:     "Pontifex Quintus Caecilius",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Expand Olive Trade Routes",
			Description: "Negotiate new trade agreements for olive oil export to Egypt",
			Status:      Vetoed,
			CreatedAt:   time.Now().Add(-60 * time.Hour),
			UpdatedAt:   time.Now().Add(-36 * time.Hour),
			Priority:    5,
			Sponsor:     "Senator Gaius Sempronius",
		},
	}

	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Prepare the insert statement
	insertStmt, err := tx.Prepare(`
		INSERT INTO todos (id, title, description, status, created_at, updated_at, priority, sponsor)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	// Insert each todo
	for _, todo := range initialTodos {
		_, err = insertStmt.Exec(
			todo.ID,
			todo.Title,
			todo.Description,
			todo.Status,
			todo.CreatedAt.Format(time.RFC3339),
			todo.UpdatedAt.Format(time.RFC3339),
			todo.Priority,
			todo.Sponsor,
		)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// GetAll returns all todos
func (r *SQLiteRepository) GetAll() []Todo {
	rows, err := r.db.Query(`
		SELECT id, title, description, status, created_at, updated_at, priority, sponsor
		FROM todos
		ORDER BY priority ASC, created_at DESC
	`)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		return []Todo{}
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

// GetByID returns a todo by its ID
func (r *SQLiteRepository) GetByID(id string) (Todo, bool) {
	row := r.db.QueryRow(`
		SELECT id, title, description, status, created_at, updated_at, priority, sponsor
		FROM todos 
		WHERE id = ?
	`, id)

	todo, err := r.scanTodo(row)
	if err != nil {
		log.Printf("Error getting todo by ID: %v", err)
		return Todo{}, false
	}

	return todo, true
}

// Create adds a new todo
func (r *SQLiteRepository) Create(todo Todo) Todo {
	// Generate ID if not provided
	if todo.ID == "" {
		todo.ID = uuid.NewString()
	}

	// Set timestamps
	now := time.Now()
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = now
	}
	todo.UpdatedAt = now

	// Default to Proposed status if not set
	if todo.Status == "" {
		todo.Status = Proposed
	}

	_, err := r.db.Exec(`
		INSERT INTO todos (id, title, description, status, created_at, updated_at, priority, sponsor)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.CreatedAt.Format(time.RFC3339),
		todo.UpdatedAt.Format(time.RFC3339),
		todo.Priority,
		todo.Sponsor,
	)

	if err != nil {
		log.Printf("Error creating todo: %v", err)
	}

	return todo
}

// Update modifies an existing todo
func (r *SQLiteRepository) Update(todo Todo) (Todo, bool) {
	// Check if the todo exists
	_, exists := r.GetByID(todo.ID)
	if !exists {
		return Todo{}, false
	}

	// Update the todo's updated_at timestamp
	todo.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		UPDATE todos
		SET title = ?, description = ?, status = ?, updated_at = ?, priority = ?, sponsor = ?
		WHERE id = ?
	`,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.UpdatedAt.Format(time.RFC3339),
		todo.Priority,
		todo.Sponsor,
		todo.ID,
	)

	if err != nil {
		log.Printf("Error updating todo: %v", err)
		return Todo{}, false
	}

	return todo, true
}

// Delete removes a todo
func (r *SQLiteRepository) Delete(id string) bool {
	result, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

// GetByStatus returns todos filtered by status
func (r *SQLiteRepository) GetByStatus(status Status) []Todo {
	rows, err := r.db.Query(`
		SELECT id, title, description, status, created_at, updated_at, priority, sponsor
		FROM todos
		WHERE status = ?
		ORDER BY priority ASC, created_at DESC
	`, status)
	if err != nil {
		log.Printf("Error getting todos by status: %v", err)
		return []Todo{}
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

// FilterByPriority returns todos with the specified priority
func (r *SQLiteRepository) FilterByPriority(priority int) []Todo {
	rows, err := r.db.Query(`
		SELECT id, title, description, status, created_at, updated_at, priority, sponsor
		FROM todos
		WHERE priority = ?
		ORDER BY created_at DESC
	`, priority)
	if err != nil {
		log.Printf("Error getting todos by priority: %v", err)
		return []Todo{}
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

// Helper function to scan multiple rows into Todo slice
func (r *SQLiteRepository) scanTodos(rows *sql.Rows) []Todo {
	var todos []Todo

	for rows.Next() {
		var (
			id           string
			title        string
			description  string
			status       string
			createdAtStr string
			updatedAtStr string
			priority     int
			sponsor      string
		)

		if err := rows.Scan(&id, &title, &description, &status, &createdAtStr, &updatedAtStr, &priority, &sponsor); err != nil {
			log.Printf("Error scanning todo row: %v", err)
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
		updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

		todos = append(todos, Todo{
			ID:          id,
			Title:       title,
			Description: description,
			Status:      Status(status),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Priority:    priority,
			Sponsor:     sponsor,
		})
	}

	return todos
}

// Helper function to scan a single row into a Todo
func (r *SQLiteRepository) scanTodo(row *sql.Row) (Todo, error) {
	var (
		id           string
		title        string
		description  string
		status       string
		createdAtStr string
		updatedAtStr string
		priority     int
		sponsor      string
	)

	if err := row.Scan(&id, &title, &description, &status, &createdAtStr, &updatedAtStr, &priority, &sponsor); err != nil {
		return Todo{}, err
	}

	createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      Status(status),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Priority:    priority,
		Sponsor:     sponsor,
	}, nil
}
