package formatter

import (
	"fmt"
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Sqlite struct{}

func (s Sqlite) FormatNode(w io.Writer, n ast.Node) {
	switch tn := n.(type) {
	case *ast.Select:
		s.formatSelect(w, tn)
	case *ast.Delete:
		s.formatDelete(w, tn)
	case *ast.Insert:
		s.formatInsert(w, tn)
	case *ast.Update:
		s.formatUpdate(w, tn)
	case *ast.CreateTable:
		s.formatCreateTable(w, tn)
	case *ast.ColumnSpec:
		s.formatColumnSpec(w, tn)
	case ast.ColumnType:
		s.formatColumnType(w, tn)
	case *ast.ColumnDefault:
		s.formatColumnDefault(w, tn)
	case ast.Nullability:
		s.formatNullability(w, tn)
	case *ast.PrimaryKey:
		// Primary keys are added to columns in the column definition during a CREATE TABLE
		return
	case *ast.TableName:
		s.formatTableName(w, tn)
	case *ast.Join:
		s.formatJoin(w, tn)
	case *ast.Alias:
		s.formatAlias(w, tn)
	case *ast.TableAlias:
		s.formatTableAlias(w, tn)
	case *ast.Identifier:
		s.formatIdentifier(w, tn)
	case *ast.Selector:
		s.formatSelector(w, tn)
	case *ast.ValuesLiteral:
		s.formatValuesLiteral(w, tn)
	case *ast.Limit:
		s.formatLimit(w, tn)
	case *ast.Lock:
		s.formatLock(w, tn)
	case *ast.Where:
		s.formatWhere(w, tn)
	case *ast.UnaryExpr:
		s.formatUnaryExpr(w, tn)
	case *ast.BinaryExpr:
		s.formatBinaryExpr(w, tn)
	case *ast.PlaceholderLiteral:
		s.formatPlaceholderLiteral(w, tn)
	case *ast.TupleLiteral:
		s.formatTupleLiteral(w, tn)
	case *ast.IntegerLiteral:
		s.formatIntegerLiteral(w, tn)
	case *ast.StringLiteral:
		s.formatStringLiteral(w, tn)
	case *ast.NullLiteral:
		s.formatNullLiteral(w, tn)
	case *ast.OrderBy:
		s.formatOrderBy(w, tn)
	case *ast.Function:
		s.formatFunction(w, tn)
	case *ast.StarLiteral:
		fmt.Fprint(w, "*")
	case *ast.Distinct:
		s.formatDistinct(w, tn)
	case *ast.OnDuplicateKey:
		s.formatOnDuplicateKey(w, tn)
	default:
		panic(fmt.Sprintf(`unexpected node: %T`, n))
	}
}

func (s Sqlite) formatSelect(w io.Writer, sl *ast.Select) {
	fmt.Fprint(w, `SELECT `)
	for i, expr := range sl.Exprs {
		s.FormatNode(w, expr)
		if i < len(sl.Exprs)-1 {
			fmt.Fprint(w, `,`)
		}
	}

	fmt.Fprint(w, ` FROM `)
	s.FormatNode(w, sl.From)

	if sl.Where != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Where)
	}
	if sl.OrderBy != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.OrderBy)
	}
	if sl.Limit != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Limit)
	}
	if sl.Lock != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Lock)
	}
}

func (s Sqlite) formatDelete(w io.Writer, d *ast.Delete) {
	fmt.Fprint(w, `DELETE FROM `)
	s.FormatNode(w, d.From)

	if d.Where != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.Where)
	}
	if d.OrderBy != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.OrderBy)
	}
	if d.Limit != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.Limit)
	}
}

func (s Sqlite) formatInsert(w io.Writer, i *ast.Insert) {
	fmt.Fprint(w, `INSERT INTO `)
	s.FormatNode(w, i.Into)
	fmt.Fprint(w, ` (`)
	formatCommaDelimited(w, s, i.Columns...)
	fmt.Fprint(w, `) VALUES `)
	formatCommaDelimited(w, s, i.Values...)
	if i.OnDuplicateKey != nil {
		s.FormatNode(w, i.OnDuplicateKey)
	}
}

func (s Sqlite) formatUpdate(w io.Writer, u *ast.Update) {
	fmt.Fprint(w, `UPDATE `)
	s.FormatNode(w, u.Table)
	fmt.Fprint(w, ` SET `)
	formatCommaDelimited(w, s, u.AssignmentList...)

	if u.Where != nil {
		fmt.Fprintf(w, ` `)
		s.FormatNode(w, u.Where)
	}

	// TODO we "support" these in the formatter, but we don't expose them to the public via the builders.
	// Add tests for these once we support them publicly.
	if u.OrderBy != nil {
		fmt.Fprint(w, ` ORDER BY `)
		s.FormatNode(w, u.OrderBy)
	}
	if u.Limit != nil {
		fmt.Fprint(w, ` LIMIT `)
		s.FormatNode(w, u.Limit)
	}
}

func (s Sqlite) formatCreateTable(w io.Writer, ct *ast.CreateTable) {
	fmt.Fprint(w, `CREATE TABLE `)
	if ct.IfNotExists {
		fmt.Fprint(w, `IF NOT EXISTS `)
	}

	s.FormatNode(w, ct.Name)

	fmt.Fprint(w, `(`)
	formatCommaDelimited(w, s, ct.Columns...)

	fmt.Fprint(w, `)`)
}

func (s Sqlite) formatColumnSpec(w io.Writer, cs *ast.ColumnSpec) {
	s.FormatNode(w, cs.Name)
	fmt.Fprint(w, ` `)
	s.FormatNode(w, cs.Type)
	if cs.Nullability != ast.NoNullability {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, cs.Nullability)
	}
	if cs.Default != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, cs.Default)
	}
	if cs.ComprisesPrimaryKey {
		fmt.Fprint(w, ` PRIMARY KEY`)
	}

	// SQLite has no concept of auto_increment
}

func (s Sqlite) formatColumnType(w io.Writer, ct ast.ColumnType) {
	switch ct.(type) {
	case ast.TinyIntColumn, ast.SmallIntColumn, ast.IntColumn, ast.BigIntColumn:
		fmt.Fprint(w, `INTEGER`)
	case ast.CharColumn, ast.VarCharColumn, ast.TextColumn:
		fmt.Fprint(w, `TEXT`)
	case ast.TinyBlobColumn, ast.BlobColumn, ast.MediumBlobColumn, ast.LongBlobColumn:
		fmt.Fprint(w, `BLOB`)
	case ast.DateTimeColumn:
		fmt.Fprint(w, `NUMERIC`)
	}
}

func (s Sqlite) formatColumnDefault(w io.Writer, cd *ast.ColumnDefault) {
	fmt.Fprint(w, `DEFAULT `)
	s.FormatNode(w, cd.Value)
}

func (s Sqlite) formatNullability(w io.Writer, n ast.Nullability) {
	switch n {
	case ast.NotNull:
		fmt.Fprint(w, `NOT NULL`)
	case ast.Null:
		fmt.Fprint(w, `NULL`)
	}
}

func (s Sqlite) formatFunction(w io.Writer, f *ast.Function) {
	fmt.Fprint(w, f.Name)
	fmt.Fprint(w, `(`)
	for _, arg := range f.Args {
		s.FormatNode(w, arg)
	}
	fmt.Fprint(w, `)`)
}

func (s Sqlite) formatIntegerLiteral(w io.Writer, l *ast.IntegerLiteral) {
	fmt.Fprintf(w, `%d`, l.Value)
}

func (s Sqlite) formatStringLiteral(w io.Writer, l *ast.StringLiteral) {
	fmt.Fprintf(w, `'%s'`, l.Value)
}

func (s Sqlite) formatNullLiteral(w io.Writer, _ *ast.NullLiteral) {
	fmt.Fprint(w, `NULL`)
}

func (s Sqlite) formatTupleLiteral(w io.Writer, t *ast.TupleLiteral) {
	fmt.Fprint(w, `(`)
	for i, val := range t.Values {
		s.FormatNode(w, val)
		if i < len(t.Values)-1 {
			fmt.Fprint(w, `,`)
		}
	}
	fmt.Fprint(w, `)`)
}

func (s Sqlite) formatPlaceholderLiteral(w io.Writer, _ *ast.PlaceholderLiteral) {
	fmt.Fprint(w, `?`)
}

func (s Sqlite) formatWhere(w io.Writer, wh *ast.Where) {
	fmt.Fprint(w, `WHERE `)
	s.FormatNode(w, wh.Expr)
}

func (s Sqlite) formatOrderBy(w io.Writer, o *ast.OrderBy) {
	fmt.Fprint(w, `ORDER BY `)
	for i, ord := range o.Orders {
		s.FormatNode(w, ord.Expr)
		switch ord.Direction {
		case ast.OrderAsc:
			fmt.Fprint(w, ` ASC`)
		case ast.OrderDesc:
			fmt.Fprint(w, ` DESC`)
		}

		if i < len(o.Orders)-1 {
			fmt.Fprint(w, `,`)
		}
	}
}

func (s Sqlite) formatLock(w io.Writer, l *ast.Lock) {
	// SQLite doesn't support locking. It doesn't support multiple concurrent transactions, so every select is more or less equivalent to MySQL's FOR UPDATE.
}

func (s Sqlite) formatLimit(w io.Writer, l *ast.Limit) {
	fmt.Fprint(w, `LIMIT `)
	if l.Offset != nil {
		s.FormatNode(w, l.Offset)
		fmt.Fprint(w, `, `)
	}
	s.FormatNode(w, l.Count)
}

func (s Sqlite) formatIdentifier(w io.Writer, c *ast.Identifier) {
	fmt.Fprint(w, c.Name)
}

func (s Sqlite) formatSelector(w io.Writer, sel *ast.Selector) {
	s.FormatNode(w, sel.SelectFrom)
	fmt.Fprint(w, ".")
	s.FormatNode(w, sel.FieldName)
}

func (s Sqlite) formatValuesLiteral(w io.Writer, vl *ast.ValuesLiteral) {
	s.FormatNode(w, &ast.Selector{
		SelectFrom: ast.NewIdentifier("excluded"),
		FieldName:  vl.Target,
	})
}

func (s Sqlite) formatTableName(w io.Writer, tn *ast.TableName) {
	s.FormatNode(w, tn.Identifier)
}

func (s Sqlite) formatJoin(w io.Writer, j *ast.Join) {
	s.FormatNode(w, j.Left)

	switch j.Kind {
	case ast.JoinKindInner:
		fmt.Fprint(w, ` INNER JOIN `)
	case ast.JoinKindLeft:
		fmt.Fprint(w, ` LEFT JOIN `)
	default:
		panic(`unexpected join kind`)
	}

	s.FormatNode(w, j.Right)
	fmt.Fprint(w, ` ON `)
	s.FormatNode(w, j.On)
}

func (s Sqlite) formatAlias(w io.Writer, a *ast.Alias) {
	s.FormatNode(w, a.ForExpr)
	fmt.Fprint(w, ` AS `)
	s.FormatNode(w, a.As)
}

func (s Sqlite) formatTableAlias(w io.Writer, a *ast.TableAlias) {
	s.FormatNode(w, a.Alias)
}

func (s Sqlite) formatUnaryExpr(w io.Writer, un *ast.UnaryExpr) {
	// For postfix operators, we format the operand before the operator.
	if un.Op.IsPost() {
		s.FormatNode(w, un.Operand)
	}

	switch un.Op {
	case ast.UnaryIsNotNull:
		fmt.Fprint(w, " IS NOT NULL")
	case ast.UnaryIsNull:
		fmt.Fprint(w, " IS NULL")
	default:
		panic(fmt.Sprintf("unexpected ast.UnaryExprOperator: %#v", un.Op))
	}

	// For prefix operators, we format the operand after the operator.
	if !un.Op.IsPost() {
		s.FormatNode(w, un.Operand)
	}
}

func (s Sqlite) formatBinaryExpr(w io.Writer, bin *ast.BinaryExpr) {
	s.FormatNode(w, bin.Left)

	switch bin.Op {
	case ast.BinaryEquals:
		fmt.Fprint(w, ` = `)
	case ast.BinaryNotEquals:
		fmt.Fprint(w, ` != `)
	case ast.BinaryGreater:
		fmt.Fprint(w, ` > `)
	case ast.BinaryGreaterOrEqual:
		fmt.Fprint(w, ` >= `)
	case ast.BinaryLess:
		fmt.Fprint(w, ` < `)
	case ast.BinaryLessOrEqual:
		fmt.Fprint(w, ` <= `)
	case ast.BinaryIn:
		fmt.Fprint(w, ` IN `)
	case ast.BinaryAnd:
		fmt.Fprint(w, ` AND `)
	case ast.BinaryOr:
		fmt.Fprint(w, ` OR `)
	default:
		panic(fmt.Sprintf(`unsupported binary operation: %v`, bin.Op))
	}

	s.FormatNode(w, bin.Right)
}

func (s Sqlite) formatDistinct(w io.Writer, d *ast.Distinct) {
	fmt.Fprint(w, `DISTINCT `)
	formatCommaDelimited(w, s, d.Exprs...)
}

func (s Sqlite) formatOnDuplicateKey(w io.Writer, odk *ast.OnDuplicateKey) {
	fmt.Fprint(w, `ON CONFLICT (`)
	formatCommaDelimited(w, s, odk.KeyIdents...)
	fmt.Fprint(w, `) DO UPDATE SET `)
	formatCommaDelimited(w, s, odk.Updates...)
}
