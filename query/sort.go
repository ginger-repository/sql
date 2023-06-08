package query

import (
	"github.com/ginger-core/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (q *query) Sort() (*gorm.DB, errors.Error) {
	dbQ := q.getDB()

	var err errors.Error
	if sub := q.GetSub(); sub != nil {
		dbQ, err = New(sub, dbQ).Sort()
		if err != nil {
			return nil, err
		}
	}

	if sorts := q.GetSorts(); sorts != nil {
		sort := sorts.NextSort()
		for sort != nil {
			dbQ = dbQ.Order(clause.OrderByColumn{
				Column: clause.Column{
					Name: sort.GetSortBy(),
				},
				Desc: !sort.IsAsc(),
			})
			sort = sorts.NextSort()
		}
	}
	if sort := q.GetSort(); sort != nil {
		dbQ = dbQ.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: sort.GetSortBy(),
			},
			Desc: !sort.IsAsc(),
		})
	}
	return dbQ, nil
}
