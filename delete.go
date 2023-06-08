package sql

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
	"github.com/ginger-repository/sql/query"
)

func (repo *repo) Delete(q _query.Query) errors.Error {
	model := repo.getModel(q)
	if model == nil {
		return ModelNotDefinedErr
	}

	db, err := repo.newDB(q, model)
	if err != nil {
		return err
	}

	qry := query.New(q, db)

	err = qry.Filter()
	if err != nil {
		return err
	}

	db = db.Callback().Delete().Execute(db)
	if db.Error != nil {
		return errors.Internal(db.Error)
	}
	if db.RowsAffected == 0 {
		return errors.NotFound().
			WithTrace("Delete().Execute.RowsAffected")
	}
	return nil
}
