package sql

import (
	"github.com/ginger-core/errors"
	"github.com/ginger-core/query"
	"gorm.io/gorm"
)

func (repo *repo) getModel(q query.Query) any {
	if q, ok := q.(query.ModelQuery); ok {
		if model := q.GetModel(); model != nil {
			return model
		}
	}
	return nil
}

func (repo *repo) getModels(q query.Query) any {
	if q, ok := q.(query.ModelsQuery); ok {
		if models := q.GetModels(); models != nil {
			return models
		}
	}
	return nil
}

func (repo *repo) newDB(q query.Query, model any) (*gorm.DB, errors.Error) {
	db := repo.getDB(q)

	if model == nil {
		if db != nil {
			return db, nil
		}
		return repo.DB, nil
	}

	if db == nil {
		db = repo.Model(model)
	}
	if db.Error != nil {
		return nil, errors.Internal(db.Error)
	}

	if ctx := q.GetContext(); ctx != nil {
		db = db.WithContext(q.GetContext())
	}

	return db, nil
}
