package ast

type Where struct {
	Expr Expr
}

func (w *Where) AcceptVisitor(fn func(Node) bool) {
	if w == nil {
		return
	}
	if fn(w) {
		w.Expr.AcceptVisitor(fn)
	}
}
