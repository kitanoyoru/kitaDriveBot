package sqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/kitanoyoru/kitaDriveBot/libs/transactor"
)

func NewTransactor(db *sqlx.DB) transactor.Transactor {
	return &sqlTransactor{db}
}

type sqlTransactor struct {
	db *sqlx.DB
}

func (t *sqlTransactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	err = tFunc(InjectTx(ctx, tx))
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
