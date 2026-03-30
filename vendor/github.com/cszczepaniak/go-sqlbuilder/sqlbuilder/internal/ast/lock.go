package ast

type LockKind int

const (
	NoLock LockKind = iota
	SharedLock
	ForUpdateLock
)

type Lock struct {
	Node
	Kind LockKind
}

func (l *Lock) AcceptVisitor(fn func(Node) bool) {
	fn(l)
}
