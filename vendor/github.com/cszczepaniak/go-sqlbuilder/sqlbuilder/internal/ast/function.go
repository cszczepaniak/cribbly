package ast

type Function struct {
	Expr
	Name string
	Args []Expr
}

func NewFunction(name string, args ...IntoExpr) *Function {
	return &Function{
		Name: name,
		Args: IntoExprs(args...),
	}
}

func (f *Function) AcceptVisitor(fn func(Node) bool) {
	if fn(f) {
		for _, a := range f.Args {
			a.AcceptVisitor(fn)
		}
	}
}
