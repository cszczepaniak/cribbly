package ast

type Update struct {
	Table          TableExpr
	AssignmentList []Expr
	Where          *Where
	OrderBy        *OrderBy
	Limit          *Limit
}

func NewUpdate(table IntoTableExpr) *Update {
	return &Update{
		Table: table.IntoTableExpr(),
	}
}

func (u *Update) AcceptVisitor(fn func(n Node) bool) {
	if fn(u) {
		u.Table.AcceptVisitor(fn)
		for _, expr := range u.AssignmentList {
			expr.AcceptVisitor(fn)
		}

		if u.Where != nil {
			u.Where.AcceptVisitor(fn)
		}
		if u.OrderBy != nil {
			u.OrderBy.AcceptVisitor(fn)
		}
		if u.Limit != nil {
			u.Limit.AcceptVisitor(fn)
		}
	}
}

func (u *Update) AddAssignments(exprs ...IntoExpr) {
	for _, expr := range exprs {
		u.AssignmentList = append(u.AssignmentList, expr.IntoExpr())
	}
}

func (u *Update) WithWhere(expr IntoExpr) *Update {
	e := expr.IntoExpr()
	if e == nil {
		return u
	}

	u.Where = &Where{
		Expr: e,
	}
	return u
}

func (u *Update) WithOrders(os ...Order) *Update {
	if u.OrderBy == nil {
		u.OrderBy = &OrderBy{
			Orders: os,
		}
		return u
	}
	u.OrderBy.Orders = append(u.OrderBy.Orders, os...)
	return u
}

func (u *Update) WithLimit(offset, count IntoExpr) *Update {
	if count == nil {
		return u
	}
	u.Limit = &Limit{
		Offset: offset.IntoExpr(),
		Count:  count.IntoExpr(),
	}
	return u
}
