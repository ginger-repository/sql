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
		return err.WithTrace("newDB")
	}

	update := q.GetUpdate()
	if update != nil {
		sets := update.GetSets()
		du := make(clause.Set, len(sets))
		for i, s := range sets {
			du[i] = clause.Assignment{
				Column: clause.Column{
					Name: s.Key,
				},
				Value: s.Value,
			}
		}
		db = db.
			Clauses(clause.OnConflict{
				DoUpdates: du,
			})
	} else {
		db = db.
			Clauses(clause.OnConflict{
				UpdateAll: true,
			})
	}

	db = db.Create(entity)
	if db.Error != nil {
		return errors.Internal(db.Error)
	}

	return nil
}
