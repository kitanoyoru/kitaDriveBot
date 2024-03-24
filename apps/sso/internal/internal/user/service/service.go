package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user"
	"github.com/kitanoyoru/kitaDriveBot/libs/hasher"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
	txLib "github.com/kitanoyoru/kitaDriveBot/libs/transactor"
)

func New(storage user.Storage, hasher hasher.Hasher, transactor txLib.Transactor, logger *logger.Logger) user.Service {
	return &service{
		storage:    storage,
		hasher:     hasher,
		transactor: transactor,
		logger:     logger,
	}
}

type service struct {
	storage    user.Storage
	hasher     hasher.Hasher
	transactor txLib.Transactor
	logger     *logger.Logger
}

func (s *service) ListUsers(ctx context.Context, req user.ListUsersRequest) ([]user.User, error) {
	return s.storage.ListUsers(ctx, user.WithIDs(req.IDs))
}

func (s *service) CreateUser(ctx context.Context, req user.CreateUserRequest) (user.User, error) {
	now := time.Now()

	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return user.User{}, err

	}

	u := user.User{
		ID:             uuid.NewString(),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	var newUser user.User
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		newUser, err = s.storage.CreateUser(ctx, u)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return user.User{}, err
	}

	return newUser, nil
}

func (s *service) UpdateUser(ctx context.Context, req user.UpdateUserRequest) (user.User, error) {
	u, err := s.storage.GetUser(ctx, req.ID)
	if err != nil {
		return user.User{}, err
	}

	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Email = req.Email

	var updatedUser user.User
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		updatedUser, err = s.storage.UpdateUser(ctx, u)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return user.User{}, err
	}

	return updatedUser, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		return s.storage.DeleteUser(ctx, id)
	})
}
