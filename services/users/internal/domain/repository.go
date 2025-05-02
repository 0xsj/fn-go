package domain

import "context"

type User struct {
    ID       string
    Email    string
    Name     string
    Password string
    Role     string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) (*User, error)
    Delete(ctx context.Context, id string) error
}