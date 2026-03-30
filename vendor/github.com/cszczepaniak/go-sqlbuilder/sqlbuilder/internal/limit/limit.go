package limit

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Limiter interface {
	Limit() (string, error)
}

// LimitBuilder is used to construct SQL limits. This builder is meant to be embedded in any other builder which can
// have a limit (e.g. Select)
//
// Note: normally, this name shouldn't stutter with the package name, but we may want to embed several things called
// "Builder" in other builders, so we have to disambiguate.
type LimitBuilder[T any] struct {
	parent T

	limit *int
}

func NewBuilder[T any](parent T) *LimitBuilder[T] {
	return &LimitBuilder[T]{
		parent: parent,
	}
}

func (b *LimitBuilder[T]) Limit(limit int) T {
	b.limit = &limit
	return b.parent
}

func (b *LimitBuilder[T]) OffsetAndLimit() (ast.IntoExpr, ast.IntoExpr) {
	if b.limit == nil {
		return nil, nil
	}
	return ast.None(), ast.NewIntegerLiteral(*b.limit)
}
