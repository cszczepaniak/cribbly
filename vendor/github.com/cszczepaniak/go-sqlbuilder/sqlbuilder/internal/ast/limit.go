package ast

type Limit struct {
	Expr
	Offset Expr
	Count  Expr
}

func NewLimit(offset, count IntoExpr) *Limit {
	return &Limit{
		Offset: offset.IntoExpr(),
		Count:  count.IntoExpr(),
	}
}

func (l *Limit) AcceptVisitor(fn func(Node) bool) {
	if l == nil {
		return
	}
	if fn(l) {
		if l.Offset != nil {
			l.Offset.AcceptVisitor(fn)
		}
		l.Count.AcceptVisitor(fn)
	}
}
