package ast

type IntoTableExpr interface {
	IntoTableExpr() TableExpr
}

type TableExpr interface {
	Expr
	tableExpr()
}

type TableName struct {
	TableExpr
	*Identifier
}

func NewTableName(name string) *TableName {
	return &TableName{
		Identifier: NewIdentifier(name),
	}
}

func (t *TableName) IntoTableExpr() TableExpr {
	return t
}

func (t *TableName) AcceptVisitor(fn func(Node) bool) {
	fn(t)
}

type TupleLiteral struct {
	TableExpr

	Values []Expr
}

func NewTupleLiteral(values ...IntoExpr) *TupleLiteral {
	return &TupleLiteral{
		Values: IntoExprs(values...),
	}
}

func (t *TupleLiteral) IntoExpr() Expr {
	return t
}

func (t *TupleLiteral) AcceptVisitor(fn func(Node) bool) {
	if fn(t) {
		for _, v := range t.Values {
			v.AcceptVisitor(fn)
		}
	}
}
