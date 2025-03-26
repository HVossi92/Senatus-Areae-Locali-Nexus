package models

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for todo operations
type Repository interface {
	GetAll() []Todo
	GetByID(id string) (Todo, bool)
	Create(todo Todo) Todo
	Update(todo Todo) (Todo, bool)
	Delete(id string) bool
	GetByStatus(status Status) []Todo
	FilterByPriority(priority int) []Todo
}

// InMemoryRepository implements Repository with an in-memory data store
type InMemoryRepository struct {
	todos map[string]Todo
	mu    sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository with sample data
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		todos: make(map[string]Todo),
	}

	// Add some initial senate themed tasks
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

	for _, todo := range initialTodos {
		repo.todos[todo.ID] = todo
	}

	return repo
}

// GetAll returns all todos
func (r *InMemoryRepository) GetAll() []Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos
}

// GetByID returns a todo by its ID
func (r *InMemoryRepository) GetByID(id string) (Todo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, found := r.todos[id]
	return todo, found
}

// Create adds a new todo
func (r *InMemoryRepository) Create(todo Todo) Todo {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a new ID if not provided
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

	r.todos[todo.ID] = todo
	return todo
}

// Update modifies an existing todo
func (r *InMemoryRepository) Update(todo Todo) (Todo, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return Todo{}, false
	}

	// Preserve creation time
	todo.CreatedAt = r.todos[todo.ID].CreatedAt
	todo.UpdatedAt = time.Now()

	r.todos[todo.ID] = todo
	return todo, true
}

// Delete removes a todo
func (r *InMemoryRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return false
	}

	delete(r.todos, id)
	return true
}

// GetByStatus returns todos filtered by status
func (r *InMemoryRepository) GetByStatus(status Status) []Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Todo
	for _, todo := range r.todos {
		if todo.Status == status {
			result = append(result, todo)
		}
	}
	return result
}

// FilterByPriority returns todos with the specified priority
func (r *InMemoryRepository) FilterByPriority(priority int) []Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Todo
	for _, todo := range r.todos {
		if todo.Priority == priority {
			result = append(result, todo)
		}
	}
	return result
}
