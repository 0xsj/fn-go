// pkg/repository/repository.go
package repository

import (
	"context"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Entity interface {
	GetID() string
}

type Repository[T Entity] interface {
	FindByID(ctx context.Context, id string) (T, error)
	
	// FindAll retrieves multiple entities with filtering
	FindAll(ctx context.Context, opts ...QueryOption) ([]T, error)
	
	// Count returns the number of entities
	Count(ctx context.Context, opts ...QueryOption) (int64, error)
	
	// Create persists a new entity
	Create(ctx context.Context, entity T) error
	
	// Update modifies an existing entity
	Update(ctx context.Context, entity T) error
	
	// Delete removes an entity
	Delete(ctx context.Context, id string) error
	
	// Transaction executes operations in a transaction
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
	logger     log.Logger
	entityName string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(logger log.Logger, entityName string) BaseRepository {
	return BaseRepository{
		logger:     logger,
		entityName: entityName,
	}
}

// NotFoundError creates a not found error
func (r *BaseRepository) NotFoundError(id string) error {
	return errors.NewNotFoundError(r.entityName+" not found", nil).
		WithField("id", id).
		WithField("entity", r.entityName)
}

// ConflictError creates a conflict error
func (r *BaseRepository) ConflictError(id string) error {
	return errors.NewConflictError(r.entityName+" already exists", nil).
		WithField("id", id).
		WithField("entity", r.entityName)
}

// DatabaseError creates a database error
func (r *BaseRepository) DatabaseError(operation string, err error) error {
	return errors.NewDatabaseError("database error during "+operation, err).
		WithOperation(operation).
		WithField("entity", r.entityName)
}

// Logger returns a logger with repository context
func (r *BaseRepository) Logger(ctx context.Context, operation string, fields ...map[string]any) log.Logger {
	logger := r.logger.With("operation", operation).
		With("entity", r.entityName)
	
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			logger = logger.With(k, v)
		}
	}
	
	return logger
}