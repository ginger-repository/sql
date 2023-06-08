package sql

import (
	"github.com/ginger-core/errors"
	_query "github.com/ginger-core/query"
	"github.com/ginger-repository/sql/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repo) extendBitwiseClause(
	dest map[string]interface{}, key string, expr clause.Expr) {
	expFace := dest[key]
	if expFace == nil {
		// does not exist, assign and return
		expr.SQL = key + " " + expr.SQL
		dest[key] = expr
		return
	}
	// extend existing
	exp := expFace.(clause.Expr)
	// extend query
	exp.SQL = "(" + exp.SQL + ")" + expr.SQL
	// extend vars
	exp.Vars = append(exp.Vars, expr.Vars...)
	// reassign
	dest[key] = exp
}

func (repo *repo) update(db *gorm.DB, q _query.Update) errors.Error {
	dest, ok := db.Statement.Dest.(map[string]interface{})
	if !ok || dest == nil {
		dest = make(map[string]interface{})
	}
	for _, u := range q.GetSets() {
		dest[u.Key] = u.Value
	}
	for _, u := range q.GetIncreases() {
		dest[u.Key] = clause.Expr{
			SQL:                u.Key + " + ?",
			Vars:               []interface{}{u.Value},
			WithoutParentheses: false,
		}
	}
	for _, u := range q.GetDecreases() {
		dest[u.Key] = clause.Expr{
			SQL:                u.Key + " - ?",
			Vars:               []interface{}{u.Value},
			WithoutParentheses: false,
		}
	}
	for _, u := range q.GetAnds() {
		dest[u.Key] = clause.Expr{
			SQL:                " & ?",
			Vars:               []interface{}{u.Value},
			WithoutParentheses: false,
		}
	}
	for _, u := range q.GetOrs() {
		repo.extendBitwiseClause(dest, u.Key,
			clause.Expr{
				SQL:                " | ?",
				Vars:               []interface{}{u.Value},
				WithoutParentheses: false,
			})
	}
	for _, u := range q.GetNots() {
		repo.extendBitwiseClause(dest, u.Key,
			clause.Expr{
				SQL:                " & ~?",
				Vars:               []interface{}{u.Value},
				WithoutParentheses: false,
			})
	}
	db.Statement.Dest = dest
	return nil
}

func (repo *repo) Update(q _query.Query, update any) errors.Error {
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

	if update != nil {
		db.Statement.Dest = update
	}
	if update := q.GetUpdate(); update != nil {
		if err := repo.update(db, update); err != nil {
			return err
		}
	}

	db = db.Callback().Update().Execute(db)
	if db.Error != nil {
		return errors.Internal(db.Error)
	}

	if db.RowsAffected == 0 {
		return errors.NotFound().
			WithTrace("Update().Execute.RowsAffected=0")
	}
	return nil
}
