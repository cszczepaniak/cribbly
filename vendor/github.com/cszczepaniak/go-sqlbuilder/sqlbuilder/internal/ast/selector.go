package ast

type Selector struct {
	Expr
	SelectFrom *Identifier
	FieldName  *Identifier
}

func (s *Selector) AcceptVisitor(fn func(Node) bool) {
	fn(s)
}

func (s *Selector) IntoExpr() Expr {
	return s
}
