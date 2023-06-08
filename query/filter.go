package query

import (
	"github.com/ginger-core/errors"
	"gorm.io/gorm/clause"
)

func (q *query) Filter() errors.Error {
	dbQ := q.getDB()

	expressions, err := q.generateExpressions(q)
	if err != nil {
		return err
	}
	if len(expressions) > 0 {
		dbQ.Statement.AddClause(clause.Where{
			Exprs: expressions,
		})
	}
	return nil
}
