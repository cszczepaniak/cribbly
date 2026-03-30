package ast

type IntoExpr interface {
	IntoExpr() Expr
}

func IntoExprs(intos ...IntoExpr) []Expr {
	res := make([]Expr, 0, len(intos))
	for _, i := range intos {
		res = append(res, i.IntoExpr())
	}
	return res
}

func None() intoNone {
	return intoNone{}
}

type intoNone struct{}

func (intoNone) IntoExpr() Expr {
	return nil
}

type Expr interface {
	Node
	expr()
}

type UnaryExprOperator int

const (
	UnaryIsNull    UnaryExprOperator = iota
	UnaryIsNotNull UnaryExprOperator = iota
)

func (op UnaryExprOperator) IsPost() bool {
	switch op {
	case UnaryIsNotNull:
		return true
	case UnaryIsNull:
		return true
	default:
		return false
	}
}

type UnaryExpr struct {
	Expr
	Op      UnaryExprOperator
	Operand Expr
}

func NewUnaryExpr(operand Expr, op UnaryExprOperator) *UnaryExpr {
	return &UnaryExpr{
		Op:      op,
		Operand: operand,
	}
}

func (u *UnaryExpr) IntoExpr() Expr {
	return u
}

func (u *UnaryExpr) AcceptVisitor(fn func(Node) bool) {
	if fn(u) {
		u.Operand.AcceptVisitor(fn)
	}
}

type BinaryExprOperator int

const (
	BinaryEquals BinaryExprOperator = iota
	BinaryNotEquals
	BinaryGreater
	BinaryGreaterOrEqual
	BinaryLess
	BinaryLessOrEqual
	BinaryIn
	BinaryAnd
	BinaryOr
)

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    BinaryExprOperator
	Right Expr
}

func NewBinaryExpr(left IntoExpr, op BinaryExprOperator, right IntoExpr) *BinaryExpr {
	return &BinaryExpr{
		Left:  left.IntoExpr(),
		Op:    op,
		Right: right.IntoExpr(),
	}
}

func (b *BinaryExpr) IntoExpr() Expr {
	return b
}

func (b *BinaryExpr) AcceptVisitor(fn func(Node) bool) {
	if fn(b) {
		b.Left.AcceptVisitor(fn)
		b.Right.AcceptVisitor(fn)
	}
}
