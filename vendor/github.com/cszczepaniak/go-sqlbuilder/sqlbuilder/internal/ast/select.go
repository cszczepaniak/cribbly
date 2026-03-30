package ast

type Select struct {
	From    TableExpr
	Exprs   []Expr
	Where   *Where
	Limit   *Limit
	OrderBy *OrderBy
	Lock    *Lock
}

func NewSelect(from TableExpr, exprs ...IntoExpr) *Select {
	es := make([]Expr, 0, len(exprs))
	for _, e := range exprs {
		es = append(es, e.IntoExpr())
	}
	return &Select{
		From:  from,
		Exprs: es,
	}
}

func (s *Select) AcceptVisitor(fn func(n Node) bool) {
	if fn(s) {
		s.From.AcceptVisitor(fn)
		for _, exp := range s.Exprs {
			exp.AcceptVisitor(fn)
		}
		s.Where.AcceptVisitor(fn)
		s.Limit.AcceptVisitor(fn)
		s.OrderBy.AcceptVisitor(fn)
		s.Lock.AcceptVisitor(fn)
	}
}

func (s *Select) WithWhere(expr IntoExpr) *Select {
	e := expr.IntoExpr()
	if e == nil {
		return s
	}

	s.Where = &Where{
		Expr: expr.IntoExpr(),
	}
	return s
}

func (s *Select) WithOrders(os ...Order) *Select {
	if s.OrderBy == nil {
		s.OrderBy = &OrderBy{
			Orders: os,
		}
		return s
	}
	s.OrderBy.Orders = append(s.OrderBy.Orders, os...)
	return s
}

func (s *Select) WithLimit(offset, count IntoExpr) *Select {
	if count == nil {
		return s
	}
	s.Limit = &Limit{
		Offset: offset.IntoExpr(),
		Count:  count.IntoExpr(),
	}
	return s
}

func (s *Select) WithLock(k LockKind) *Select {
	s.Lock = &Lock{Kind: k}
	return s
}
