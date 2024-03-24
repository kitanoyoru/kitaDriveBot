package user

import "context"

type Storage interface {
	ListUsers(ctx context.Context, filters ...ListUsersFilter) ([]User, error)
	GetUser(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, req User) (User, error)
	UpdateUser(ctx context.Context, req User) (User, error)
	DeleteUser(ctx context.Context, id string) error
}

type ListUsersFilters struct {
	IDs []string
}

type ListUsersFilter func(*ListUsersFilters)

func WithIDs(ids []string) ListUsersFilter {
	return func(f *ListUsersFilters) {
		f.IDs = ids
	}
}
