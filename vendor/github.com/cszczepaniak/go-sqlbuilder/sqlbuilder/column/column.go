package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type ColumnExpressionBuilder struct {
	ident *ast.Identifier
}

func (b *ColumnExpressionBuilder) IntoExpr() ast.Expr {
	return b.ident.IntoExpr()
}

func Named(name string) *ColumnExpressionBuilder {
	return &ColumnExpressionBuilder{
		ident: ast.NewIdentifier(name),
	}
}

func (b *ColumnExpressionBuilder) QualifiedBy(qualifier string) ast.IntoExpr {
	return &ast.Selector{
		SelectFrom: ast.NewIdentifier(qualifier),
		FieldName:  b.ident,
	}
}
