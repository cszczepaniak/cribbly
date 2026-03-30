package table

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

// JoinBuilder is used to express an in-progress join.
type JoinBuilder struct {
	kind       ast.JoinKind
	joiningTo  *TableBuilder
	toBeJoined ast.IntoTableExpr
}

// On completes the join. It specifies a condition used to join the two table expressions.
func (jb *JoinBuilder) On(expr ast.Expr) *TableBuilder {
	return newTableBuilder(ast.NewJoin(
		jb.kind,
		jb.joiningTo.IntoTableExpr(),
		jb.toBeJoined.IntoTableExpr(),
		expr,
	))
}

// OnEqualExpressions completes the join using equality of the given expressions as the join
// condition.
func (jb *JoinBuilder) OnEqualExpressions(left, right ast.IntoExpr) *TableBuilder {
	return jb.On(ast.NewBinaryExpr(
		left,
		ast.BinaryEquals,
		right,
	))
}

// OnEqualColumns completes the join using equality of the given column names as the join
// condition.
func (jb *JoinBuilder) OnEqualColumns(left, right string) *TableBuilder {
	return jb.OnEqualExpressions(
		ast.NewIdentifier(left),
		ast.NewIdentifier(right),
	)
}
