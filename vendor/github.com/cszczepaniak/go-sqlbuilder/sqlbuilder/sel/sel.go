package sel

import (
	"context"
	"database/sql"
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/condition"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/limit"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type Builder struct {
	tableExpr ast.IntoTableExpr
	forUpdate bool
	orderBy   *filter.Order

	exprs []ast.IntoExpr

	*condition.ConditionBuilder[*Builder]
	*limit.LimitBuilder[*Builder]

	formatter Formatter
}

func NewBuilder(f Formatter, tableExpr ast.IntoTableExpr) *Builder {
	b := &Builder{
		tableExpr: tableExpr,
		formatter: f,
	}

	b.ConditionBuilder = condition.NewBuilder(b)
	b.LimitBuilder = limit.NewBuilder(b)
	return b
}

func (b *Builder) Columns(colNames ...string) *Builder {
	for _, f := range colNames {
		b.exprs = append(b.exprs, ast.NewIdentifier(f))
	}
	return b
}

func (b *Builder) Expressions(exprs ...ast.IntoExpr) *Builder {
	b.exprs = append(b.exprs, exprs...)
	return b
}

func (b *Builder) ForUpdate() *Builder {
	b.forUpdate = true
	return b
}

func (b *Builder) OrderBy(o filter.Order) *Builder {
	b.orderBy = &o
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	n := ast.NewSelect(b.tableExpr.IntoTableExpr(), b.exprs...)

	n.WithWhere(b.ConditionBuilder)

	offset, limit := b.LimitBuilder.OffsetAndLimit()
	n.WithLimit(offset, limit)

	if b.orderBy != nil {
		n.WithOrders(ast.NewOrder(ast.NewIdentifier(b.orderBy.Column), b.orderBy.Direction.ToASTDirection()))
	}

	if b.forUpdate {
		n.WithLock(ast.ForUpdateLock)
	}

	sb := &strings.Builder{}
	b.formatter.FormatNode(sb, n)

	return statement.Statement{
		Stmt: sb.String(),
		Args: ast.GetArgs(n),
	}, nil
}

func (b *Builder) Query(q dispatch.Queryer) (*sql.Rows, error) {
	return dispatch.Query(b, q)
}

func (b *Builder) QueryContext(ctx context.Context, q dispatch.QueryCtxer) (*sql.Rows, error) {
	return dispatch.QueryContext(ctx, b, q)
}

func (b *Builder) QueryRow(q dispatch.RowQueryer) (*sql.Row, error) {
	return dispatch.QueryRow(b, q)
}

func (b *Builder) QueryRowContext(ctx context.Context, q dispatch.RowQueryCtxer) (*sql.Row, error) {
	return dispatch.QueryRowContext(ctx, b, q)
}
