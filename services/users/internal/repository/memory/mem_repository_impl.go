package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/0xsj/fn-go/services/users/internal/domain"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrEmailExists  = errors.New("email already exists")
)

type UserRepository struct {
    users  map[string]*domain.User
    emails map[string]string // email -> id
    mutex  sync.RWMutex
    nextID int
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        users:  make(map[string]*domain.User),
        emails: make(map[string]string),
        nextID: 1,
    }
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    // Check if email already exists
    if _, exists := r.emails[user.Email]; exists {
        return nil, ErrEmailExists
    }

    // Create a new ID if not provided
    if user.ID == "" {
        user.ID = generateID(r.nextID)
        r.nextID++
    }

    // Store the user
    newUser := &domain.User{
        ID:       user.ID,
        Email:    user.Email,
        Name:     user.Name,
        Password: user.Password,
        Role:     user.Role,
    }

    r.users[newUser.ID] = newUser
    r.emails[newUser.Email] = newUser.ID

    return newUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, ErrUserNotFound
    }

    return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    existingUser, exists := r.users[user.ID]
    if !exists {
        return nil, ErrUserNotFound
    }

    // Check if email is being changed and if it's already in use
    if user.Email != existingUser.Email {
        if _, emailExists := r.emails[user.Email]; emailExists {
            return nil, ErrEmailExists
        }
        // Remove old email mapping
        delete(r.emails, existingUser.Email)
        // Add new email mapping
        r.emails[user.Email] = user.ID
    }

    // Update user
    updatedUser := &domain.User{
        ID:       user.ID,
        Email:    user.Email,
        Name:     user.Name,
        Password: existingUser.Password, // Don't update password here
        Role:     user.Role,
    }

    r.users[user.ID] = updatedUser

    return updatedUser, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    user, exists := r.users[id]
    if !exists {
        return ErrUserNotFound
    }

    // Remove email mapping
    delete(r.emails, user.Email)
    // Remove user
    delete(r.users, id)

    return nil
}

func generateID(id int) string {
    return fmt.Sprintf("user-%d", id)
}