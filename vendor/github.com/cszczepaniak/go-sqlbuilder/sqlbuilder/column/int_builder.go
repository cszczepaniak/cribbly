package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type anyInteger interface {
	int8 | int16 | int32 | int64
}

type integerColumnBuilder[T anyInteger, U columnTyper] struct {
	*baseColumnBuilder[T, U]
	autoIncrement bool
}

func newIntegerColumnBuilder[T anyInteger, U columnTyper](name string, parent U) *integerColumnBuilder[T, U] {
	return &integerColumnBuilder[T, U]{
		baseColumnBuilder: newBaseColumnBuilder[T](name, parent),
	}
}

func (b *integerColumnBuilder[T, U]) AutoIncrement() U {
	b.autoIncrement = true
	return b.parent
}

func (b *integerColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := b.baseColumnBuilder.Build().SetAutoIncrement(b.autoIncrement)

	if b.defaultVal != nil {
		cs.WithDefault(ast.NewIntegerLiteral(int(*b.defaultVal)))
	}

	return cs
}
