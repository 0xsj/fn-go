// pkg/repository/repository.go
package repository

import (
	"context"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
)

// Entity represents the minimal interface for entities stored in repositories
type Entity interface {
    GetID() string
}

// Repository defines the basic operations that any repository should support
type Repository[T Entity] interface {
    // FindByID retrieves an entity by its ID
    FindByID(ctx context.Context, id string) (T, error)
    
    // FindAll retrieves all entities, optionally filtered
    FindAll(ctx context.Context, opts ...QueryOption) ([]T, error)
    
    // Count returns the number of entities, optionally filtered
    Count(ctx context.Context, opts ...QueryOption) (int64, error)
    
    // Create persists a new entity
    Create(ctx context.Context, entity T) error
    
    // Update modifies an existing entity
    Update(ctx context.Context, entity T) error
    
    // Delete removes an entity by ID
    Delete(ctx context.Context, id string) error
    
    // Transaction executes the provided function within a transaction
    Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BaseRepository provides common functionality for repositories
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
    msg := r.entityName + " not found"
    return errors.NewNotFoundError(msg, nil).WithField("id", id)
}

// ConflictError creates a conflict error
func (r *BaseRepository) ConflictError(id string) error {
    msg := r.entityName + " already exists"
    return errors.NewConflictError(msg, nil).WithField("id", id)
}

// DatabaseError creates a database error
func (r *BaseRepository) DatabaseError(operation string, err error) error {
    msg := "Database error during " + operation
    return errors.NewDatabaseError(msg, err).
        WithOperation(operation).
        WithField("entity", r.entityName)
}

// LogOperation logs a repository operation
func (r *BaseRepository) LogOperation(ctx context.Context, operation string, fields ...map[string]interface{}) log.Logger {
    opLogger := r.logger.With("operation", operation).With("entity", r.entityName)
    
    if len(fields) > 0 && fields[0] != nil {
        for k, v := range fields[0] {
            opLogger = opLogger.With(k, v)
        }
    }
    
    return opLogger
}