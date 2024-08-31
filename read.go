package sql

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
	"github.com/ginger-repository/sql/query"
)

func (repo *repo) Count(q _query.Query) (uint64, errors.Error) {
	models := repo.getModels(q)
	if models == nil {
		return 0, ModelNotDefinedErr
	}

	db, err := repo.newDB(q, models)
	if err != nil {
		return 0, err
	}

	qry := query.New(q, db)

	err = qry.Filter()
	if err != nil {
		return 0, err
	}

	var count int64
	db = db.Count(&count)
	if db.Error != nil {
		return 0, errors.Internal(db.Error)
	}

	return uint64(count), nil
}

func (repo *repo) List(q _query.Query) (any, errors.Error) {
	models := repo.getModels(q)
	db, err := repo.newDB(q, models)
	if err != nil {
		return nil, err
	}

	qry := query.New(q, db)

	err = qry.Filter()
	if err != nil {
		return 0, err
	}

	db, err = qry.Sort()
	if err != nil {
		return nil, err
	}

	pagination := qry.GetPagination()
	if pagination != nil {
		db = db.
			Offset((pagination.GetPage() - 1) * pagination.GetSize()).
			Limit(pagination.GetSize())
	}

	db = db.Find(models)
	if db.Error != nil {
		return nil, errors.Internal(db.Error)
	}

	return models, nil
}

func (repo *repo) Get(q _query.Query) (any, errors.Error) {
	model := repo.getModel(q)
	db, err := repo.newDB(q, model)
	if err != nil {
		return nil, err
	}

	qry := query.New(q, db)

	err = qry.Filter()
	if err != nil {
		return nil, err
	}

	db, err = qry.Sort()
	if err != nil {
		return nil, err
	}

	db = db.Limit(1).Find(model)
	if db.Error != nil {
		return nil, errors.Internal(db.Error)
	}

	if db.RowsAffected == 0 {
		return nil, errors.NotFound(errors.DefaultNotFoundError)
	}

	return model, nil
}
