package handler

import (
	"context"

	"github.com/0xsj/fn-go/pkg/common/logging"
	"github.com/0xsj/fn-go/pkg/proto/users"
	"github.com/0xsj/fn-go/services/users/internal/config"
	"github.com/0xsj/fn-go/services/users/internal/domain"
	"github.com/0xsj/fn-go/services/users/internal/repository/memory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type UserServiceHandler struct {
// 	users.UnimplementedUserServiceServer
// 	cfg *config.Config
// 	logger logging.Logger
// }

// func NewUserServiceHandler(cfg *config.Config, logger logging.Logger) *UserServiceHandler {
// 	return &UserServiceHandler{
// 		cfg: cfg,
// 		logger: logger,
// 	}
// }

// func (h *UserServiceHandler) GetUser(ctx context.Context, req *users.GetUserRequest) (*users.GetUserResponse, error) {
// 	h.logger.Info("GetUser request received", logging.F("id", req.Id))
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "user ID is required")
// 	}

// 	user := &users.User{
// 		Id: req.Id,
// 		Email: "user@example.com",
// 		Name: "Example User",
// 		Role: "user",
// 	}

// 	return &users.GetUserResponse{
// 		User: user,
// 	}, nil
// }


type UserServiceHandler struct {
    users.UnimplementedUserServiceServer
    cfg    *config.Config
    logger logging.Logger
    repo   domain.UserRepository
}

func NewUserServiceHandler(cfg *config.Config, logger logging.Logger) *UserServiceHandler {
    handlerLogger := logger.With(logging.F("component", "user-handler"))
    repo := memory.NewUserRepository()
    
    return &UserServiceHandler{
        cfg:    cfg,
        logger: handlerLogger,
        repo:   repo,
    }
}

func (h *UserServiceHandler) GetUser(ctx context.Context, req *users.GetUserRequest) (*users.GetUserResponse, error) {
    h.logger.Info("GetUser request received", logging.F("id", req.Id))
    
    // Validate request
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "user ID is required")
    }
    
    // Get user from repository
    user, err := h.repo.GetByID(ctx, req.Id)
    if err != nil {
        if err == memory.ErrUserNotFound {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        h.logger.Error("Failed to get user", logging.F("error", err))
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    // Convert domain user to proto user
    protoUser := &users.User{
        Id:    user.ID,
        Email: user.Email,
        Name:  user.Name,
        Role:  user.Role,
    }
    
    return &users.GetUserResponse{
        User: protoUser,
    }, nil
}

func (h *UserServiceHandler) CreateUser(ctx context.Context, req *users.CreateUserRequest) (*users.CreateUserResponse, error) {
    h.logger.Info("CreateUser request received", logging.F("email", req.Email))
    
    // Validate request
    if req.Email == "" || req.Name == "" || req.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "email, name, and password are required")
    }
    
    // Create user in repository
    user := &domain.User{
        Email:    req.Email,
        Name:     req.Name,
        Password: req.Password, // In a real app, you'd hash this
        Role:     req.Role,
    }
    
    createdUser, err := h.repo.Create(ctx, user)
    if err != nil {
        if err == memory.ErrEmailExists {
            return nil, status.Error(codes.AlreadyExists, "email already exists")
        }
        h.logger.Error("Failed to create user", logging.F("error", err))
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    // Convert domain user to proto user
    protoUser := &users.User{
        Id:    createdUser.ID,
        Email: createdUser.Email,
        Name:  createdUser.Name,
        Role:  createdUser.Role,
    }
    
    return &users.CreateUserResponse{
        User: protoUser,
    }, nil
}

func (h *UserServiceHandler) UpdateUser(ctx context.Context, req *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
    h.logger.Info("UpdateUser request received", logging.F("id", req.Id))
    
    // Validate request
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "user ID is required")
    }
    
    // Update user in repository
    user := &domain.User{
        ID:    req.Id,
        Email: req.Email,
        Name:  req.Name,
        Role:  req.Role,
    }
    
    updatedUser, err := h.repo.Update(ctx, user)
    if err != nil {
        if err == memory.ErrUserNotFound {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        if err == memory.ErrEmailExists {
            return nil, status.Error(codes.AlreadyExists, "email already exists")
        }
        h.logger.Error("Failed to update user", logging.F("error", err))
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    // Convert domain user to proto user
    protoUser := &users.User{
        Id:    updatedUser.ID,
        Email: updatedUser.Email,
        Name:  updatedUser.Name,
        Role:  updatedUser.Role,
    }
    
    return &users.UpdateUserResponse{
        User: protoUser,
    }, nil
}

func (h *UserServiceHandler) DeleteUser(ctx context.Context, req *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
    h.logger.Info("DeleteUser request received", logging.F("id", req.Id))
    
    // Validate request
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "user ID is required")
    }
    
    // Delete user in repository
    err := h.repo.Delete(ctx, req.Id)
    if err != nil {
        if err == memory.ErrUserNotFound {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        h.logger.Error("Failed to delete user", logging.F("error", err))
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    return &users.DeleteUserResponse{
        Success: true,
    }, nil
}