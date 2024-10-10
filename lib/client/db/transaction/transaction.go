package transaction

import (
	"context"
	"github.com/creamcx/help/lib/client/db"
	"github.com/creamcx/help/lib/client/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactions
}

func NewTransactionManager(db db.Transactions) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			err = errors.Wrap(err, "recover transaction")
		}
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrap(errRollback, "rollback transaction")
			}
			return
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "commit transaction")
			}
		}

	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "transaction")
	}
	return err
}

func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
