// services/auth-service/internal/repository/mysql/role_permission_repository.go
package repository

import (
	"context"
	"database/sql"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
)

type RolePermissionRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewRolePermissionRepository(db *sql.DB, logger log.Logger) repository.RolePermissionRepository {
	return &RolePermissionRepository{
		db:     db,
		logger: logger.WithLayer("mysql-role-permission-repository"),
	}
}

func (r *RolePermissionRepository) AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id, created_at)
		VALUES (?, ?, NOW())
		ON DUPLICATE KEY UPDATE created_at = created_at
	`
	
	_, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to assign permission to role"),
			"assign_permission_to_role",
		)
	}
	
	return nil
}

func (r *RolePermissionRepository) RevokePermissionFromRole(ctx context.Context, roleID string, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`
	
	result, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to revoke permission from role"),
			"revoke_permission_from_role",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"revoke_permission_from_role",
		)
	}
	
	if affected == 0 {
		return domain.NewPermissionNotFoundError("role_permission")
	}
	
	return nil
}

func (r *RolePermissionRepository) GetRolesByPermission(ctx context.Context, permissionID string) ([]string, error) {
	query := `SELECT role_id FROM role_permissions WHERE permission_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, permissionID)
	if err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get roles by permission"),
			"get_roles_by_permission",
		)
	}
	defer rows.Close()
	
	var roleIDs []string
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, domain.WithOperation(
				domain.Wrap(err, "failed to scan role ID"),
				"get_roles_by_permission",
			)
		}
		roleIDs = append(roleIDs, roleID)
	}
	
	if err := rows.Err(); err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "error iterating role rows"),
			"get_roles_by_permission",
		)
	}
	
	return roleIDs, nil
}