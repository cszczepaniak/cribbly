package update

import (
	"context"
	"database/sql"
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/condition"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type fieldAndArg struct {
	field string
	arg   any
}

type Builder struct {
	f Formatter

	table  ast.IntoTableExpr
	fields []fieldAndArg

	*condition.ConditionBuilder[*Builder]
}

func NewBuilder(f Formatter, table ast.IntoTableExpr) *Builder {
	b := &Builder{
		table: table,
		f:     f,
	}

	b.ConditionBuilder = condition.NewBuilder(b)
	return b
}

func (b *Builder) SetFieldTo(field string, val any) *Builder {
	b.fields = append(b.fields, fieldAndArg{
		field: field,
		arg:   val,
	})
	return b
}

func (b *Builder) SetFieldToNull(field string) *Builder {
	b.fields = append(b.fields, fieldAndArg{
		field: field,
		arg:   nil,
	})
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	u := ast.NewUpdate(b.table)

	exprs := make([]ast.IntoExpr, 0, len(b.fields))
	for _, field := range b.fields {
		var rhs ast.IntoExpr
		if field.arg != nil {
			rhs = ast.NewPlaceholderLiteral(field.arg)
		} else {
			rhs = ast.NewNullLiteral()
		}

		b := ast.NewBinaryExpr(
			ast.NewIdentifier(field.field),
			ast.BinaryEquals,
			rhs,
		)
		exprs = append(exprs, b)
	}

	u.AddAssignments(exprs...)
	u.WithWhere(b.ConditionBuilder)

	sb := strings.Builder{}
	b.f.FormatNode(&sb, u)

	return statement.Statement{
		Stmt: sb.String(),
		Args: ast.GetArgs(u),
	}, nil
}

func (b *Builder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *Builder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
