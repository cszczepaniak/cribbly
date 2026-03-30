package ast

type JoinKind int

const (
	JoinKindLeft JoinKind = iota
	JoinKindInner
)

type Join struct {
	TableExpr

	Kind  JoinKind
	Left  TableExpr
	Right TableExpr
	On    Expr
}

func NewJoin(
	kind JoinKind,
	left TableExpr,
	right TableExpr,
	on Expr,
) *Join {
	return &Join{
		Kind:  kind,
		Left:  left,
		Right: right,
		On:    on,
	}
}

func (t *Join) IntoExpr() Expr {
	return t
}

func (t *Join) IntoTableExpr() TableExpr {
	return t
}

func (t *Join) AcceptVisitor(fn func(Node) bool) {
	fn(t)
}
