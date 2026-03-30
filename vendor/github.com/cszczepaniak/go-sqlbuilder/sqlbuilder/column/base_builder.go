package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type columnTyper interface {
	columnType() ast.ColumnType
}

type baseColumnBuilder[T any, U columnTyper] struct {
	name        string
	defaultVal  *T
	defaultNull bool
	nullable    *bool
	primaryKey  bool

	parent U
}

func newBaseColumnBuilder[T any, U columnTyper](name string, parent U) *baseColumnBuilder[T, U] {
	return &baseColumnBuilder[T, U]{
		name:   name,
		parent: parent,
	}
}

func (b *baseColumnBuilder[T, U]) Default(val T) U {
	b.defaultVal = &val
	return b.parent
}

func (b *baseColumnBuilder[T, U]) DefaultNull() U {
	b.defaultNull = true
	b.defaultVal = nil
	return b.parent
}

func (b *baseColumnBuilder[T, U]) Null() U {
	tr := true
	b.nullable = &tr
	return b.parent
}

func (b *baseColumnBuilder[T, U]) NotNull() U {
	f := false
	b.nullable = &f
	return b.parent
}

func (b *baseColumnBuilder[T, U]) PrimaryKey() U {
	b.primaryKey = true
	return b.parent
}

func (b *baseColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := ast.NewColumnSpec(b.name, b.parent.columnType()).
		WithNullabilityFromBool(b.nullable)

	if b.defaultNull {
		cs.WithDefault(ast.NewNullLiteral())
	}
	cs.SetPrimaryKey(b.primaryKey)

	return cs
}
