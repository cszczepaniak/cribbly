package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type stringColumnBuilder[U columnTyper] struct {
	*baseColumnBuilder[string, U]
	size int
}

func newStringColumnBuilder[U columnTyper](name string, size int, parent U) *stringColumnBuilder[U] {
	return &stringColumnBuilder[U]{
		baseColumnBuilder: newBaseColumnBuilder[string](name, parent),
		size:              size,
	}
}

func (b *stringColumnBuilder[U]) Build() *ast.ColumnSpec {
	cs := b.baseColumnBuilder.Build()

	if b.defaultVal != nil {
		cs.WithDefault(ast.NewStringLiteral(*b.defaultVal))
	}

	return cs
}
