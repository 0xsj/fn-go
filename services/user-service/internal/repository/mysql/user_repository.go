// services/user-service/internal/repository/mysql/user_repository.go
package mysql

import (
	"database/sql"

	"github.com/0xsj/fn-go/pkg/common/log"
)

// UserRepository implements repository.UserRepository using MySQL
type UserRepository struct {
	db     *sql.DB
	logger log.Logger
}

// NewUserRepository creates a new MySQL user repository
// func NewUserRepository(db *sql.DB, logger log.Logger) repository.UserRepository {
// 	return &UserRepository{
// 		db:     db,
// 		logger: logger.WithLayer("mysql-user-repository"),
// 	}
// }
