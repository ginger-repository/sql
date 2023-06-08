package query

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Query interface {
	_query.Query
	getDB() *gorm.DB

	GetModel() any
	GetModels() any
	FindModel() (any, errors.Error)

	Filter() errors.Error
	Sort() (*gorm.DB, errors.Error)
}

type query struct {
	_query.Query

	DB *gorm.DB
}

func New(source _query.Query, db *gorm.DB) Query {
	return &query{
		Query: source,
		DB:    db,
	}
}

func (q *query) SetDB(db any) {
	q.DB = db.(*gorm.DB)
}

func (q *query) GetDB() any {
	return q.getDB()
}

func (q *query) getDB() *gorm.DB {
	if q.DB == nil {
		return nil
	}
	return q.DB
}

func (q *query) generateExpressions(
	_q _query.Query) ([]clause.Expression, errors.Error) {
	exp := make([]clause.Expression, 0)
	if filter := _q.GetFilter(); filter != nil {
		if matchs := filter.GetMatchs(); len(matchs) > 0 {
			matchsExps := make([]clause.Expression, 0)
			for _, m := range matchs {
				if m.HasCustomHandle() {
					result, err := m.HandleCustom(q)
					if err != nil {
						return nil, err.WithTrace("m.HandleCustom")
					}
					if result != nil {
						exps := result.GetData().([]clause.Expression)
						exp = append(exp, exps...)
					}
				} else {
					e := q.getMatchExpression(m)
					if e != nil {
						matchsExps = append(matchsExps, e)
					}
				}
			}
			if len(matchsExps) > 0 {
				exp = append(exp, clause.And(matchsExps...))
			}
		}
		for _, op := range filter.GetOperations() {
			f := op.GetFilter()
			if f == nil {
				continue
			}
			if matchs := filter.GetMatchs(); len(matchs) > 0 {
				filters := make([]clause.Expression, 0)
				for _, m := range matchs {
					if m.HasCustomHandle() {
						result, err := m.HandleCustom(q)
						if err != nil {
							return nil, err.WithTrace("m.HandleCustom")
						}
						if result != nil {
							exps := result.GetData().([]clause.Expression)
							if exps != nil {
								filters = exps
							}
						}
					} else {
						e := q.getMatchExpression(m)
						if e != nil {
							filters = append(filters, e)
						}
					}
				}
				if len(exp) > 0 {
					switch op.GetOperation() {
					case _query.OperationAnd:
						clauses := append(exp, filters...)
						exp = []clause.Expression{
							clause.And(clauses...),
						}
						// exp = append(exp, clause.And(filters...))
					case _query.OperationOr:
						clauses := append(exp, filters...)
						exp = []clause.Expression{
							clause.Or(clauses...),
						}
						// exp = append(exp, clause.Or(filters...))
					}
				}
			} else {
				exps, err := q.generateExpressions(f)
				if err != nil {
					return nil, err
				}
				if len(exps) > 0 {
					switch op.GetOperation() {
					case _query.OperationAnd:
						clauses := append(exp, exps...)
						exp = []clause.Expression{
							clause.And(clauses...),
						}
						// exp = append(exp, clause.And(exps...))
					case _query.OperationOr:
						clauses := append(exp, exps...)
						exp = []clause.Expression{
							clause.Or(clauses...),
						}
						// exp = append(exp, clause.Or(exps...))
					}
				}
			}
		}
		if sub := filter.GetSub(); sub != nil {
			exps, err := q.generateExpressions(sub)
			if err != nil {
				return nil, err
			}
			if len(exps) > 0 {
				if len(exp) > 0 {
					exp = []clause.Expression{
						clause.And(exp...),
					}
				}
				exp = append(exp, exps...)
			}
			// if len(exp) > 0 {
			// 	return exp, nil
			// }
		}
	}
	return exp, nil
}
