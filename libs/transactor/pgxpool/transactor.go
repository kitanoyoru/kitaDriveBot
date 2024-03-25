package pgxpool

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/kitanoyoru/kitaDriveBot/libs/transactor"
)

func NewTransactor(db *pgxpool.Pool) transactor.Transactor {
	return &sqlTransactor{db}
}

type sqlTransactor struct {
	db *pgxpool.Pool
}

func (t *sqlTransactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	err = tFunc(InjectTx(ctx, tx))
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return errors.Wrap(err, "failed to rollback transaction")
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
