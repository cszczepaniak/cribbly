package ast

type TableAlias struct {
	TableExpr
	*Alias
}

func (a *TableAlias) IntoTableExpr() TableExpr {
	return a
}

func (a *TableAlias) AcceptVisitor(fn func(Node) bool) {
	fn(a)
}
