package ast

type OrderDirection int

const (
	OrderAsc OrderDirection = iota
	OrderDesc
)

type Order struct {
	Expr      Expr
	Direction OrderDirection
}

func NewOrder(expr IntoExpr, dir OrderDirection) Order {
	return Order{
		Expr:      expr.IntoExpr(),
		Direction: dir,
	}
}

type OrderBy struct {
	Orders []Order
}

func (o *OrderBy) AcceptVisitor(fn func(Node) bool) {
	fn(o)
}
