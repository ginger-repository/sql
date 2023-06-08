package sql

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
	"github.com/go-sql-driver/mysql"
)

func (repo *repo) Create(q _query.Query, entity any) errors.Error {
	model := repo.getModel(q)
	if model == nil {
		return ModelNotDefinedErr
	}

	db, err := repo.newDB(q, model)
	if err != nil {
		return err
	}

	db = db.Create(entity)
	if db.Error != nil {
		switch err := db.Error.(type) {
		case *mysql.MySQLError:
			switch err.Number {
			case 1062:
				return errors.Duplicate(err)
			}
		}
		return errors.Internal(db.Error)
	}

	return nil
}
