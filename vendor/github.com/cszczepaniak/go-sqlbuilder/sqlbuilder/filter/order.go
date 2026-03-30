package filter

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Direction int

const (
	Ascending Direction = iota
	Descending
)

func (d Direction) ToASTDirection() ast.OrderDirection {
	switch d {
	case Ascending:
		return ast.OrderAsc
	case Descending:
		return ast.OrderDesc
	}
	panic(`unreachable`)
}

type Order struct {
	Column    string
	Direction Direction
}

func OrderDesc(field string) Order {
	return Order{
		Column:    field,
		Direction: Descending,
	}
}

func OrderAsc(field string) Order {
	return Order{
		Column:    field,
		Direction: Ascending,
	}
}
