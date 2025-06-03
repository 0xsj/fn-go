// services/auth-service/internal/repository/mysql/permission_repository.go
package repository

import (
	"context"
	"database/sql"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
)

type PermissionRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewPermissionRepository(db *sql.DB, logger log.Logger) repository.PermissionRepository {
	return &PermissionRepository{
		db:     db,
		logger: logger.WithLayer("mysql-permission-repository"),
	}
}

func (r *PermissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	query := `
		INSERT INTO permissions (id, name, description, resource, action, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.Resource,
		permission.Action,
		permission.CreatedAt,
		permission.UpdatedAt,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to create permission in database"),
			"create_permission",
		)
	}
	
	return nil
}

func (r *PermissionRepository) GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {
	query := `
		SELECT id, name, description, resource, action, created_at, updated_at
		FROM permissions
		WHERE id = ?
	`
	
	permission := &models.Permission{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewPermissionNotFoundError(id)
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get permission from database"),
			"get_permission_by_id",
		)
	}
	
	return permission, nil
}

func (r *PermissionRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	query := `
		SELECT id, name, description, resource, action, created_at, updated_at
		FROM permissions
		WHERE name = ?
	`
	
	permission := &models.Permission{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewPermissionNotFoundError(name)
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get permission by name from database"),
			"get_permission_by_name",
		)
	}
	
	return permission, nil
}

func (r *PermissionRepository) ListPermissions(ctx context.Context) ([]*models.Permission, error) {
	query := `
		SELECT id, name, description, resource, action, created_at, updated_at
		FROM permissions
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to list permissions from database"),
			"list_permissions",
		)
	}
	defer rows.Close()
	
	var permissions []*models.Permission
	for rows.Next() {
		permission := &models.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, domain.WithOperation(
				domain.Wrap(err, "failed to scan permission row"),
				"list_permissions",
			)
		}
		permissions = append(permissions, permission)
	}
	
	if err := rows.Err(); err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "error iterating permission rows"),
			"list_permissions",
		)
	}
	
	return permissions, nil
}

func (r *PermissionRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	query := `
		UPDATE permissions 
		SET name = ?, description = ?, resource = ?, action = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(
		ctx,
		query,
		permission.Name,
		permission.Description,
		permission.Resource,
		permission.Action,
		permission.UpdatedAt,
		permission.ID,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to update permission in database"),
			"update_permission",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"update_permission",
		)
	}
	
	if affected == 0 {
		return domain.NewPermissionNotFoundError(permission.ID)
	}
	
	return nil
}

func (r *PermissionRepository) DeletePermission(ctx context.Context, id string) error {
	query := `DELETE FROM permissions WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to delete permission from database"),
			"delete_permission",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"delete_permission",
		)
	}
	
	if affected == 0 {
		return domain.NewPermissionNotFoundError(id)
	}
	
	return nil
}

func (r *PermissionRepository) GetPermissionsByRole(ctx context.Context, roleID string) ([]*models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.name
	`
	
	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get permissions by role from database"),
			"get_permissions_by_role",
		)
	}
	defer rows.Close()
	
	var permissions []*models.Permission
	for rows.Next() {
		permission := &models.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, domain.WithOperation(
				domain.Wrap(err, "failed to scan permission row"),
				"get_permissions_by_role",
			)
		}
		permissions = append(permissions, permission)
	}
	
	if err := rows.Err(); err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "error iterating permission rows"),
			"get_permissions_by_role",
		)
	}
	
	return permissions, nil
}