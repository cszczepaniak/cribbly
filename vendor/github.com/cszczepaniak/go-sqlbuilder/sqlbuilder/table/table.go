package table

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type TableBuilder struct {
	tableExpr ast.IntoTableExpr
}

func newTableBuilder(tableExpr ast.IntoTableExpr) *TableBuilder {
	return &TableBuilder{
		tableExpr: tableExpr,
	}
}

func Named(name string) *TableBuilder {
	return newTableBuilder(ast.NewTableName(name))
}

func (tb *TableBuilder) IntoTableExpr() ast.TableExpr {
	return tb.tableExpr.IntoTableExpr()
}

func (tb *TableBuilder) As(alias string) *TableBuilder {
	tb.tableExpr = &ast.TableAlias{
		Alias: &ast.Alias{
			ForExpr: tb.tableExpr.IntoTableExpr(),
			As:      ast.NewIdentifier(alias),
		},
	}
	return tb
}

func (tb *TableBuilder) LeftJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindLeft,
		joiningTo:  tb,
		toBeJoined: tableExpr,
	}
}

func (tb *TableBuilder) InnerJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindInner,
		joiningTo:  tb,
		toBeJoined: tableExpr,
	}
}
