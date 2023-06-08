package query

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
)

func (q *query) GetModel() any {
	return GetModel(q)
}

func (q *query) GetModels() any {
	return GetModels(q)
}

func (q *query) FindModel() (any, errors.Error) {
	return FindModel(q)
}

func GetModel(q _query.Query) any {
	if q, ok := q.(_query.ModelQuery); ok {
		if model := q.GetModel(); model != nil {
			return model
		}
	}
	return nil
}

func GetModels(q _query.Query) any {
	if q, ok := q.(_query.ModelsQuery); ok {
		if models := q.GetModels(); models != nil {
			return models
		}
	}
	return nil
}

func FindModel(q _query.Query) (any, errors.Error) {
	if model := GetModel(q); model != nil {
		return model, nil
	}
	if models := GetModels(q); models != nil {
		return models, nil
	}
	return nil, ModelNotDefinedErr
}
