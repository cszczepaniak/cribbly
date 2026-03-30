package ast

type Alias struct {
	Expr
	ForExpr Expr
	As      *Identifier
}

func (a *Alias) IntoExpr() Expr {
	return a
}

func (a *Alias) AcceptVisitor(fn func(Node) bool) {
	fn(a)
}
