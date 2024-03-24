package user

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	ListUsers(ctx context.Context, req ListUsersRequest) ([]User, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (User, error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) (User, error)
	DeleteUser(ctx context.Context, id string) error
}

type ListUsersRequest struct {
	IDs []string
}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateUserRequest struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}
