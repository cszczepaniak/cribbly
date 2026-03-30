package ast

type Delete struct {
	From    TableExpr
	Where   *Where
	Limit   *Limit
	OrderBy *OrderBy
}

func NewDelete(from TableExpr) *Delete {
	return &Delete{
		From: from,
	}
}

func (s *Delete) AcceptVisitor(fn func(n Node) bool) {
	if fn(s) {
		s.From.AcceptVisitor(fn)
		s.Where.AcceptVisitor(fn)
		s.Limit.AcceptVisitor(fn)
		s.OrderBy.AcceptVisitor(fn)
	}
}

func (d *Delete) WithWhere(expr IntoExpr) *Delete {
	e := expr.IntoExpr()
	if e == nil {
		return d
	}

	d.Where = &Where{
		Expr: expr.IntoExpr(),
	}
	return d
}

func (d *Delete) WithOrders(os ...Order) *Delete {
	if d.OrderBy == nil {
		d.OrderBy = &OrderBy{
			Orders: os,
		}
		return d
	}
	d.OrderBy.Orders = append(d.OrderBy.Orders, os...)
	return d
}

func (d *Delete) WithLimit(offset, count IntoExpr) *Delete {
	if count == nil {
		return d
	}
	d.Limit = &Limit{
		Offset: offset.IntoExpr(),
		Count:  count.IntoExpr(),
	}
	return d
}
