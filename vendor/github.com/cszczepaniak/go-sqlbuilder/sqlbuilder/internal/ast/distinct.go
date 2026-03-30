package ast

type Distinct struct {
	Expr
	Exprs []Expr
}

func NewDistinct(exprs ...IntoExpr) *Distinct {
	return &Distinct{
		Exprs: IntoExprs(exprs...),
	}
}

func (d *Distinct) AcceptVisitor(fn func(Node) bool) {
	if fn(d) {
		for _, expr := range d.Exprs {
			fn(expr)
		}
	}
}

func (d *Distinct) IntoExpr() Expr {
	return d
}
