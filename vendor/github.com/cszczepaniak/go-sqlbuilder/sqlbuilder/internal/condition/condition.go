package condition

import (
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Conditioner interface {
	Condition(f filter.Filter) (string, error)
}

// ConditionBuilder is used to construct SQL conditions (typically a WHERE clause). This builder is meant to be embedded
// in any other builder which can have a condition (e.g. Select)
//
// Note: normally, this name shouldn't stutter with the package name, but we may want to embed several things called
// "Builder" in other builders, so we have to disambiguate.
type ConditionBuilder[T any] struct {
	parent T

	f filter.Filter
}

func NewBuilder[T any](parent T) *ConditionBuilder[T] {
	return &ConditionBuilder[T]{
		parent: parent,
	}
}

func (b *ConditionBuilder[T]) Where(f filter.Filter) T {
	b.f = f
	return b.parent
}

func (b *ConditionBuilder[T]) WhereAll(f ...filter.Filter) T {
	return b.Where(filter.All(f...))
}

func (b *ConditionBuilder[T]) WhereAny(f ...filter.Filter) T {
	return b.Where(filter.Any(f...))
}

func (b *ConditionBuilder[T]) IntoExpr() ast.Expr {
	if b.f == nil {
		return nil
	}
	return b.f.IntoExpr()
}
