package ast

type PrimaryKey struct {
	Columns []*Identifier
}

func NewPrimaryKey() *PrimaryKey {
	return &PrimaryKey{}
}

func (pk *PrimaryKey) AcceptVisitor(fn func(n Node) bool) {
	if fn(pk) {
		for _, col := range pk.Columns {
			fn(col)
		}
	}
}

func (pk *PrimaryKey) AddColumn(name string) {
	pk.Columns = append(pk.Columns, NewIdentifier(name))
}
