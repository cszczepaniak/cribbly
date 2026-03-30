package ast

type ColumnDefault struct {
	Value Expr
}

func newColumnDefault(val Expr) *ColumnDefault {
	return &ColumnDefault{
		Value: val,
	}
}

func (c *ColumnDefault) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}
