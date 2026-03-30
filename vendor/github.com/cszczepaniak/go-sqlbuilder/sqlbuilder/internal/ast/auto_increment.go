package ast

type AutoIncrement struct{}

func (a *AutoIncrement) AcceptVisitor(fn func(n Node) bool) {
	fn(a)
}
