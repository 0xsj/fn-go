// pkg/repository/memory/repository.go
package memory

import (
	"context"
	"sync"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/repository"
)

// MemoryRepository provides an in-memory implementation of Repository
type MemoryRepository[T repository.Entity] struct {
	repository.BaseRepository
	entities map[string]T
	mu       sync.RWMutex
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository[T repository.Entity](logger log.Logger, entityName string) *MemoryRepository[T] {
	return &MemoryRepository[T]{
		BaseRepository: repository.NewBaseRepository(logger, entityName),
		entities:       make(map[string]T),
	}
}

// FindByID retrieves an entity by ID
func (r *MemoryRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logger := r.LogOperation(ctx, "FindByID", map[string]interface{}{"id": id})
	logger.Debug("Looking up entity by ID")

	entity, exists := r.entities[id]
	if !exists {
		logger.Debug("Entity not found")
		var zero T
		return zero, r.NotFoundError(id)
	}

	logger.Debug("Entity found")
	return entity, nil
}

// FindAll retrieves all entities matching filters
func (r *MemoryRepository[T]) FindAll(ctx context.Context, opts ...repository.QueryOption) ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logger := r.LogOperation(ctx, "FindAll")
	options := repository.ApplyOptions(opts...)
	
	logger.With("filters", options.Filters).
		With("limit", options.Limit).
		With("offset", options.Offset).
		Debug("Finding entities with options")

	var result []T
	for _, entity := range r.entities {
		// Apply filters
		include := true
		for field, value := range options.Filters {
			// This is a simplistic approach; in a real implementation,
			// you'd use reflection or a more sophisticated way to check fields
			if field == "id" && entity.GetID() != value.(string) {
				include = false
				break
			}
		}

		if include {
			result = append(result, entity)
		}
	}

	// Apply pagination
	if options.Limit > 0 && options.Offset >= 0 && options.Offset < len(result) {
		end := options.Offset + options.Limit
		if end > len(result) {
			end = len(result)
		}
		result = result[options.Offset:end]
	}

	logger.With("count", len(result)).Debug("Found entities")
	return result, nil
}

// Count returns the number of entities
func (r *MemoryRepository[T]) Count(ctx context.Context, opts ...repository.QueryOption) (int64, error) {
	logger := r.LogOperation(ctx, "Count")
	
	entities, err := r.FindAll(ctx, opts...)
	if err != nil {
		logger.With("error", err.Error()).Error("Failed to count entities")
		return 0, err
	}
	
	count := int64(len(entities))
	logger.With("count", count).Debug("Counted entities")
	return count, nil
}

// Create adds a new entity
func (r *MemoryRepository[T]) Create(ctx context.Context, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.GetID()
	logger := r.LogOperation(ctx, "Create", map[string]interface{}{"id": id})
	logger.Debug("Creating new entity")

	if _, exists := r.entities[id]; exists {
		logger.Warn("Entity already exists")
		return r.ConflictError(id)
	}

	r.entities[id] = entity
	logger.Info("Entity created successfully")
	return nil
}

// Update modifies an existing entity
func (r *MemoryRepository[T]) Update(ctx context.Context, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.GetID()
	logger := r.LogOperation(ctx, "Update", map[string]interface{}{"id": id})
	logger.Debug("Updating entity")

	if _, exists := r.entities[id]; !exists {
		logger.Warn("Entity not found")
		return r.NotFoundError(id)
	}

	r.entities[id] = entity
	logger.Info("Entity updated successfully")
	return nil
}

// Delete removes an entity
func (r *MemoryRepository[T]) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger := r.LogOperation(ctx, "Delete", map[string]interface{}{"id": id})
	logger.Debug("Deleting entity")

	if _, exists := r.entities[id]; !exists {
		logger.Warn("Entity not found")
		return r.NotFoundError(id)
	}

	delete(r.entities, id)
	logger.Info("Entity deleted successfully")
	return nil
}

// Transaction executes the provided function (no actual transaction in memory)
func (r *MemoryRepository[T]) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	logger := r.LogOperation(ctx, "Transaction")
	logger.Debug("Executing transaction function")
	
	// No actual transaction support in memory repository
	err := fn(ctx)
	
	if err != nil {
		logger.With("error", err.Error()).Error("Transaction failed")
	} else {
		logger.Debug("Transaction completed successfully")
	}
	
	return err
}