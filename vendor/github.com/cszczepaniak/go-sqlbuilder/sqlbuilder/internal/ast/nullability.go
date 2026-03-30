package ast

type Nullability int

const (
	NoNullability Nullability = iota
	NotNull
	Null
)

func (n Nullability) AcceptVisitor(fn func(n Node) bool) {
	fn(n)
}
