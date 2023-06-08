package query

import (
	_query "github.com/ginger-core/query"
	"gorm.io/gorm/clause"
)

func GetMatchExpression(m *_query.Match) clause.Expression {
	if m == nil {
		return nil
	}
	switch m.Operator {
	case _query.None:
		return clause.Expr{SQL: m.Key, Vars: []any{m.Value}}
	case _query.Equal:
		return clause.Expr{SQL: m.Key + "=?", Vars: []any{m.Value}}
	case _query.NotEqual:
		return clause.Expr{SQL: m.Key + "!=?", Vars: []any{m.Value}}
	case _query.Is:
		switch m.Value {
		case nil:
			return clause.Expr{SQL: m.Key + " IS NULL"}
		case "1", "true", "TRUE":
			m.Value = true
		case "0", "false", "FALSE":
			m.Value = false
		}
		return clause.Expr{SQL: m.Key + " IS ?", Vars: []any{m.Value}}
	case _query.IsNot:
		switch m.Value {
		case nil:
			return clause.Expr{SQL: m.Key + " IS NOT NULL"}
		case "1", "true", "TRUE":
			m.Value = true
		case "0", "false", "FALSE":
			m.Value = false
		}
		return clause.Expr{SQL: m.Key + " IS NOT ?", Vars: []any{m.Value}}
	case _query.Lower:
		return clause.Expr{SQL: m.Key + "<?", Vars: []any{m.Value}}
	case _query.LowerEqual:
		return clause.Expr{SQL: m.Key + "<=?", Vars: []any{m.Value}}
	case _query.Greater:
		return clause.Expr{SQL: m.Key + ">?", Vars: []any{m.Value}}
	case _query.GreaterEqual:
		return clause.Expr{SQL: m.Key + ">=?", Vars: []any{m.Value}}
	case _query.BitwiseIs:
		var v any = 0
		if len(m.Extra) > 0 {
			v = m.Extra[0]
		}
		return clause.Expr{SQL: m.Key + "&?>?", Vars: append([]any{m.Value}, v)}
	case _query.BitwiseIsNot:
		return clause.Expr{SQL: m.Key + "&?=0", Vars: []any{m.Value}}
	case _query.BitwiseAndEqual:
		return clause.Expr{SQL: m.Key + "&?=?", Vars: []any{m.Value, m.Value}}
	case _query.BitwiseAndNotEqual:
		return clause.Expr{SQL: m.Key + "&?!=?", Vars: []any{m.Value, m.Value}}
	case _query.In:
		return clause.Expr{SQL: m.Key + " IN ?", Vars: []any{m.Value}}
	}
	return nil
}

func (q *query) getMatchExpression(m *_query.Match) clause.Expression {
	return GetMatchExpression(m)
}
