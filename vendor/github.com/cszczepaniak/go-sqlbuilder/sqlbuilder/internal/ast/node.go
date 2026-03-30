package ast

type Node interface {
	AcceptVisitor(fn func(n Node) bool)
}
