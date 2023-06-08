package sql

import (
	"github.com/ginger-core/errors"
	"github.com/ginger-core/query"
	"gorm.io/gorm/clause"
)

func (repo *repo) Upsert(q query.Query, entity any) errors.Error {
	model := repo.getModel(q)
	if model == nil {
		return ModelNotDefinedErr
	}

	db, err := repo.newDB(q, model)
	if err != nil {
		return err
	}

	db = db.
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(entity)
	if db.Error != nil {
		return errors.Internal(db.Error)
	}

	return nil
}
