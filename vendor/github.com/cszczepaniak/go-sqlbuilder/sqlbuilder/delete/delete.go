package delete

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
	table string
	f     Formatter

	orderBy *filter.Order
	*condition.ConditionBuilder[*Builder]
	*limit.LimitBuilder[*Builder]
}

func NewBuilder(f Formatter, table string) *Builder {
	b := &Builder{
		table: table,
		f:     f,
	}
	b.ConditionBuilder = condition.NewBuilder(b)
	b.LimitBuilder = limit.NewBuilder(b)
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	target := ast.NewTableName(b.table)
	n := ast.NewDelete(target)

	n.WithWhere(b.ConditionBuilder)

	offset, limit := b.LimitBuilder.OffsetAndLimit()
	n.WithLimit(offset, limit)

	if b.orderBy != nil {
		n.WithOrders(ast.NewOrder(ast.NewIdentifier(b.orderBy.Column), b.orderBy.Direction.ToASTDirection()))
	}

	sb := &strings.Builder{}
	b.f.FormatNode(sb, n)

	return statement.Statement{
		Stmt: sb.String(),
		Args: ast.GetArgs(n),
	}, nil
}

func (b *Builder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *Builder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
