package table

import (
	"context"
	"database/sql"
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type columnBuilder interface {
	Build() *ast.ColumnSpec
}

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type CreateBuilder struct {
	f Formatter

	name              string
	columns           []columnBuilder
	createIfNotExists bool
}

func NewCreateBuilder(f Formatter, name string) *CreateBuilder {
	return &CreateBuilder{
		f:    f,
		name: name,
	}
}

func (b *CreateBuilder) IfNotExists() *CreateBuilder {
	b.createIfNotExists = true
	return b
}

func (b *CreateBuilder) Columns(cs ...columnBuilder) *CreateBuilder {
	b.columns = append(b.columns, cs...)
	return b
}

func (b *CreateBuilder) Build() (statement.Statement, error) {
	ct := ast.NewCreateTable(b.name)
	if b.createIfNotExists {
		ct.CreateIfNotExists()
	}

	for _, col := range b.columns {
		ct.AddColumn(col.Build())
	}

	sb := &strings.Builder{}
	b.f.FormatNode(sb, ct)

	return statement.Statement{
		Stmt: sb.String(),
	}, nil
}

func (b *CreateBuilder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *CreateBuilder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
