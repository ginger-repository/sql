package sql

import (
	"database/sql"

	"github.com/ginger-core/errors"
	"github.com/ginger-core/query"
	"github.com/ginger-core/repository"
	"github.com/ginger-core/repository/tx"
)

const txTag = "sqlTx"

func (repo *repo) Begin(query query.Query,
	options ...repository.Options) errors.Error {
	if query.GetTag(txTag) == true {
		return nil
	}
	db := repo.getDB(query)
	var opts []*sql.TxOptions
	if len(options) > 0 {
		opts = make([]*sql.TxOptions, 0)
		for _, option := range options {
			switch op := option.(type) {
			case *sql.TxOptions:
				opts = append(opts, op)
			case tx.IsolationLevelOption:
				opts = append(opts, &sql.TxOptions{
					Isolation: repo.getIsolationLevel(op.GetLevel()),
					ReadOnly:  op.IsReadOnly(),
				})
			default:
				errors.Internal().
					WithTrace("Begin.argsW.cast.Err")
			}
		}
	}
	tx := db.Begin(opts...)
	query.SetDB(tx)
	query.SetTag(txTag, true)
	return nil
}

func (repo *repo) End(query query.Query) errors.Error {
	db := repo.getDB(query)
	err := query.GetError()
	if err == nil {
		if ctx := query.GetContext(); ctx != nil {
			if cErr := ctx.Err(); cErr != nil {
				err = errors.New(cErr).WithTrace("End.ctx.Err")
			}
		}
	}
	if err != nil {
		db = db.Rollback()
		if db.Error != nil {
			return errors.Internal(db.Error).
				WithTrace("End.Rollback.Err")
		}
		return err.WithTrace("End.Err")
	}
	db = db.Commit()
	if db.Error != nil {
		return errors.Internal(db.Error).
			WithTrace("End.Commit.Err")
	}
	query.DeleteTag(txTag)
	return nil
}

func (repo *repo) getIsolationLevel(level tx.IsolationLevel) sql.IsolationLevel {
	return sql.IsolationLevel(level)
}
