package ast

type ColumnSpec struct {
	Expr

	Name             *Identifier
	Type             ColumnType
	Nullability      Nullability
	Default          *ColumnDefault
	AutoIncrementing *AutoIncrement

	ComprisesPrimaryKey bool
}

func NewColumnSpec(name string, typ ColumnType) *ColumnSpec {
	return &ColumnSpec{
		Name:        NewIdentifier(name),
		Type:        typ,
		Nullability: NoNullability,
	}
}

func (cs *ColumnSpec) AcceptVisitor(fn func(n Node) bool) {
	if fn(cs) {
		cs.Name.AcceptVisitor(fn)
		cs.Type.AcceptVisitor(fn)
		if cs.Nullability != NoNullability {
			cs.Nullability.AcceptVisitor(fn)
		}
		if cs.Default != nil {
			cs.Default.AcceptVisitor(fn)
		}
		if cs.AutoIncrementing != nil {
			cs.AutoIncrementing.AcceptVisitor(fn)
		}
	}
}

func (c *ColumnSpec) WithNullabilityFromBool(val *bool) *ColumnSpec {
	if val == nil {
		return c
	}

	if *val {
		c.Nullability = Null
	} else {
		c.Nullability = NotNull
	}

	return c
}

func (c *ColumnSpec) SetAutoIncrement(val bool) *ColumnSpec {
	if val {
		c.AutoIncrementing = &AutoIncrement{}
	}
	return c
}

func (c *ColumnSpec) WithDefault(val IntoExpr) *ColumnSpec {
	c.Default = newColumnDefault(val.IntoExpr())
	return c
}

func (c *ColumnSpec) SetPrimaryKey(val bool) *ColumnSpec {
	c.ComprisesPrimaryKey = val
	return c
}
