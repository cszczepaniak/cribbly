package ast

type Identifier struct {
	Expr
	Name string
}

func NewIdentifier(name string) *Identifier {
	return &Identifier{
		Name: name,
	}
}

func (c *Identifier) AcceptVisitor(fn func(Node) bool) {
	fn(c)
}

func (c *Identifier) IntoExpr() Expr {
	return c
}
