package filter

import (
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Filter interface {
	ast.IntoExpr
}

type AllFilter struct {
	Filters []Filter
}

func All(fs ...Filter) AllFilter {
	return AllFilter{
		Filters: fs,
	}
}

func (f AllFilter) IntoExpr() ast.Expr {
	return makeChainedExpr(f.Filters[0], ast.BinaryAnd, f.Filters[1:]...).IntoExpr()
}

func makeChainedExpr(left Filter, op ast.BinaryExprOperator, rest ...Filter) ast.IntoExpr {
	if len(rest) == 0 {
		return left
	}
	return ast.NewBinaryExpr(left, op, makeChainedExpr(rest[0], op, rest[1:]...))
}

type AnyFilter struct {
	Filters []Filter
}

func Any(fs ...Filter) AnyFilter {
	return AnyFilter{
		Filters: fs,
	}
}

func (f AnyFilter) IntoExpr() ast.Expr {
	return makeChainedExpr(f.Filters[0], ast.BinaryOr, f.Filters[1:]...).IntoExpr()
}

type BinOpFilter[T any] struct {
	column string
	value  T
	op     ast.BinaryExprOperator
}

func (f BinOpFilter[T]) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewIdentifier(f.column), f.op, ast.NewPlaceholderLiteral(f.value))
}

type EqualsFilter[T any] struct {
	Column string
	Value  T
}

func Equals[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryEquals,
	}
}

func NotEquals[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryNotEquals,
	}
}

type GreaterFilter[T any] struct {
	Column string
	Value  any
}

func Greater[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryGreater,
	}
}

func GreaterOrEqual[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryGreaterOrEqual,
	}
}

func Less[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryLess,
	}
}

func LessOrEqual[T any](column string, val T) BinOpFilter[T] {
	return BinOpFilter[T]{
		column: column,
		value:  val,
		op:     ast.BinaryLessOrEqual,
	}
}

type InFilter[T any] struct {
	Column string
	Values []T
}

func In[T any](column string, vals ...T) InFilter[T] {
	return InFilter[T]{
		Column: column,
		Values: vals,
	}
}

func (f InFilter[T]) IntoExpr() ast.Expr {
	exprs := make([]ast.IntoExpr, 0, len(f.Values))
	for _, val := range f.Values {
		exprs = append(exprs, ast.NewPlaceholderLiteral(val))
	}
	return ast.NewBinaryExpr(ast.NewIdentifier(f.Column), ast.BinaryIn, ast.NewTupleLiteral(exprs...))
}

type NullFilter struct {
	column string
	op     ast.UnaryExprOperator
}

func IsNull(column string) NullFilter {
	return NullFilter{column: column, op: ast.UnaryIsNull}
}

func IsNotNull(column string) NullFilter {
	return NullFilter{column: column, op: ast.UnaryIsNotNull}
}

func (f NullFilter) IntoExpr() ast.Expr {
	return ast.NewUnaryExpr(ast.NewIdentifier(f.column), f.op)
}
