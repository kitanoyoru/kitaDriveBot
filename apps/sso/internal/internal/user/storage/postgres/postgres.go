package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/user"
	"github.com/kitanoyoru/kitaDriveBot/libs/logger"
)

type User struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func New(logger *logger.Logger, db *sqlx.DB) *User {
	return &User{
		db:     db,
		logger: logger,
	}
}

func (s *User) ListUsers(ctx context.Context, filters ...user.ListUsersFilter) ([]user.User, error) {
	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"first_name",
			"last_name",
			"email",
			"hashed_password",
			"created_at",
			"updated_at",
		).
		From("users").
		OrderBy("created_at DESC")

	listFilters := &user.ListUsersFilters{}
	for _, filter := range filters {
		filter(listFilters)
	}

	if len(listFilters.IDs) > 0 {
		query = query.Where(sq.Eq{"id": listFilters.IDs})
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	var users []user.User
	for rows.Next() {
		var user user.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.HashedPassword,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *User) CreateUser(ctx context.Context, req user.User) (user.User, error) {
	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("users").
		Columns(
			"id",
			"first_name",
			"last_name",
			"email",
			"hashed_password",
			"created_at",
			"updated_at",
		).
		Values(
			req.ID,
			req.FirstName,
			req.LastName,
			req.Email,
			req.HashedPassword,
			req.CreatedAt,
			req.UpdatedAt,
		)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return user.User{}, err
	}

	_, err = s.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return user.User{}, err
	}

	return req, nil
}

func (s *User) UpdateUser(ctx context.Context, req user.User) (user.User, error) {
	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("users").
		Set("first_name", req.FirstName).
		Set("last_name", req.LastName).
		Set("email", req.Email).
		Set("updated_at", req.UpdatedAt).
		Where(sq.Eq{"id": req.ID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return user.User{}, err
	}

	_, err = s.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return user.User{}, err
	}

	return req, nil
}

func (s *User) DeleteUser(ctx context.Context, id string) error {
	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete("users").
		Where(sq.Eq{"id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
