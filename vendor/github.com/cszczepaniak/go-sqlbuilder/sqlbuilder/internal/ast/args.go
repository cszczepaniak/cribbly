package ast

func GetArgs(n Node) []any {
	var args []any

	n.AcceptVisitor(func(n Node) bool {
		ph, ok := n.(*PlaceholderLiteral)
		if !ok {
			// Keep traversing the tree
			return true
		}
		args = append(args, ph.For)
		return false
	})

	return args
}
