// pkg/repository/sql/repository.go
package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/db"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/repository"
)

// SQLRepository provides a SQL implementation of Repository
type SQLRepository[T repository.Entity] struct {
	repository.BaseRepository
	db        db.DB
	tableName string
	
	// Mapping functions
	scanRow    func(db.Row) (T, error)
	scanRows   func(db.Rows) ([]T, error)
	toColumns  func(T) ([]string, []interface{})
	fromEntity func(T) map[string]interface{}
}

// SQLRepositoryOption defines options for creating a SQLRepository
type SQLRepositoryOption[T repository.Entity] func(*SQLRepository[T])

// NewSQLRepository creates a new SQL repository
func NewSQLRepository[T repository.Entity](
	logger log.Logger,
	database db.DB,
	entityName string,
	tableName string,
	options ...SQLRepositoryOption[T],
) *SQLRepository[T] {
	repo := &SQLRepository[T]{
		BaseRepository: repository.NewBaseRepository(logger, entityName),
		db:             database,
		tableName:      tableName,
	}
	
	for _, option := range options {
		option(repo)
	}
	
	return repo
}

// WithScanRow sets the function to scan a row into an entity
func WithScanRow[T repository.Entity](fn func(db.Row) (T, error)) SQLRepositoryOption[T] {
	return func(repo *SQLRepository[T]) {
		repo.scanRow = fn
	}
}

// WithScanRows sets the function to scan rows into entities
func WithScanRows[T repository.Entity](fn func(db.Rows) ([]T, error)) SQLRepositoryOption[T] {
	return func(repo *SQLRepository[T]) {
		repo.scanRows = fn
	}
}

// WithToColumns sets the function to convert an entity to columns and values
func WithToColumns[T repository.Entity](fn func(T) ([]string, []interface{})) SQLRepositoryOption[T] {
	return func(repo *SQLRepository[T]) {
		repo.toColumns = fn
	}
}

// WithFromEntity sets the function to convert an entity to a map
func WithFromEntity[T repository.Entity](fn func(T) map[string]interface{}) SQLRepositoryOption[T] {
	return func(repo *SQLRepository[T]) {
		repo.fromEntity = fn
	}
}

// FindByID retrieves an entity by ID
func (r *SQLRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
	logger := r.LogOperation(ctx, "FindByID", map[string]interface{}{"id": id})
	logger.Debug("Looking up entity by ID")
	
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1", r.tableName)
	row := r.db.QueryRow(ctx, query, id)
	
	entity, err := r.scanRow(row)
	if err != nil {
		// Check for no rows error
		if strings.Contains(err.Error(), "no rows") {
			logger.Debug("Entity not found")
			var zero T
			return zero, r.NotFoundError(id)
		}
		logger.With("error", err.Error()).Error("Database error during lookup")
		return entity, r.DatabaseError("FindByID", err)
	}
	
	logger.Debug("Entity found")
	return entity, nil
}

// FindAll retrieves all entities matching filters
func (r *SQLRepository[T]) FindAll(ctx context.Context, opts ...repository.QueryOption) ([]T, error) {
	logger := r.LogOperation(ctx, "FindAll")
	options := repository.ApplyOptions(opts...)
	
	logger.With("filters", options.Filters).
		With("limit", options.Limit).
		With("offset", options.Offset).
		Debug("Finding entities with options")
	
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	var args []interface{}
	
	// Add WHERE clauses for filters
	if len(options.Filters) > 0 {
		whereClauses := make([]string, 0, len(options.Filters))
		for field, value := range options.Filters {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	
	// Add ORDER BY clause for sorting
	if len(options.Sort) > 0 {
		sortClauses := make([]string, 0, len(options.Sort))
		for _, sort := range options.Sort {
			direction := "ASC"
			if sort.Direction == repository.SortDescending {
				direction = "DESC"
			}
			sortClauses = append(sortClauses, fmt.Sprintf("%s %s", sort.Field, direction))
		}
		query += " ORDER BY " + strings.Join(sortClauses, ", ")
	}
	
	// Add LIMIT and OFFSET for pagination
	if options.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", options.Limit)
		if options.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", options.Offset)
		}
	}
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		logger.With("error", err.Error()).Error("Database error during query")
		return nil, r.DatabaseError("FindAll", err)
	}
	defer rows.Close()
	
	entities, err := r.scanRows(rows)
	if err != nil {
		logger.With("error", err.Error()).Error("Error scanning rows")
		return nil, r.DatabaseError("FindAll.Scan", err)
	}
	
	logger.With("count", len(entities)).Debug("Found entities")
	return entities, nil
}

// Count returns the number of entities
func (r *SQLRepository[T]) Count(ctx context.Context, opts ...repository.QueryOption) (int64, error) {
	logger := r.LogOperation(ctx, "Count")
	options := repository.ApplyOptions(opts...)
	
	logger.With("filters", options.Filters).Debug("Counting entities with options")
	
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var args []interface{}
	
	// Add WHERE clauses for filters
	if len(options.Filters) > 0 {
		whereClauses := make([]string, 0, len(options.Filters))
		for field, value := range options.Filters {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	
	var count int64
	row := r.db.QueryRow(ctx, query, args...)
	
	err := row.Scan(&count)
	if err != nil {
		logger.With("error", err.Error()).Error("Database error during count")
		return 0, r.DatabaseError("Count", err)
	}
	
	logger.With("count", count).Debug("Counted entities")
	return count, nil
}

// Create adds a new entity
func (r *SQLRepository[T]) Create(ctx context.Context, entity T) error {
	id := entity.GetID()
	logger := r.LogOperation(ctx, "Create", map[string]interface{}{"id": id})
	logger.Debug("Creating new entity")
	
	if r.toColumns == nil {
		err := fmt.Errorf("toColumns function not set")
		logger.With("error", err.Error()).Error("Configuration error")
		return err
	}
	
	columns, values := r.toColumns(entity)
	
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	
	_, err := r.db.Execute(ctx, query, values...)
	if err != nil {
		// Check for duplicate key errors
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE constraint") {
			logger.With("error", err.Error()).Warn("Entity already exists")
			return r.ConflictError(id)
		}
		
		logger.With("error", err.Error()).Error("Database error during create")
		return r.DatabaseError("Create", err)
	}
	
	logger.Info("Entity created successfully")
	return nil
}

// Update modifies an existing entity
func (r *SQLRepository[T]) Update(ctx context.Context, entity T) error {
	id := entity.GetID()
	logger := r.LogOperation(ctx, "Update", map[string]interface{}{"id": id})
	logger.Debug("Updating entity")
	
	if r.fromEntity == nil {
		err := fmt.Errorf("fromEntity function not set")
		logger.With("error", err.Error()).Error("Configuration error")
		return err
	}
	
	data := r.fromEntity(entity)
	
	setClauses := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data))
	
	for column, value := range data {
		if column != "id" {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}
	}
	
	// Add ID for WHERE clause
	args = append(args, id)
	
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		r.tableName,
		strings.Join(setClauses, ", "),
	)
	
	rowsAffected, err := r.db.Execute(ctx, query, args...)
	if err != nil {
		logger.With("error", err.Error()).Error("Database error during update")
		return r.DatabaseError("Update", err)
	}
	
	if rowsAffected == 0 {
		logger.Warn("Entity not found for update")
		return r.NotFoundError(id)
	}
	
	logger.Info("Entity updated successfully")
	return nil
}

// Delete removes an entity
func (r *SQLRepository[T]) Delete(ctx context.Context, id string) error {
	logger := r.LogOperation(ctx, "Delete", map[string]interface{}{"id": id})
	logger.Debug("Deleting entity")
	
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	
	rowsAffected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		logger.With("error", err.Error()).Error("Database error during delete")
		return r.DatabaseError("Delete", err)
	}
	
	if rowsAffected == 0 {
		logger.Warn("Entity not found for delete")
		return r.NotFoundError(id)
	}
	
	logger.Info("Entity deleted successfully")
	return nil
}

// Transaction executes the provided function within a transaction
func (r *SQLRepository[T]) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	logger := r.LogOperation(ctx, "Transaction")
	logger.Debug("Executing transaction")
	
	return r.db.WithTransaction(ctx, func(txCtx context.Context, _ db.Transaction) error {
		return fn(txCtx)
	})
}