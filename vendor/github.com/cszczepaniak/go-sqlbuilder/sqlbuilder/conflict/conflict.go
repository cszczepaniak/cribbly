package conflict

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Key struct {
	fields []string
}

func (k Key) Fields() []string {
	return k.fields
}

func NewKey(fields ...string) Key {
	return Key{
		fields: fields,
	}
}

type Behavior interface {
	ast.IntoExpr
	Field() string
}

type IgnoreBehavior struct {
	field string
}

func Ignore(field string) IgnoreBehavior {
	return IgnoreBehavior{
		field: field,
	}
}

func (b IgnoreBehavior) Field() string { return b.field }

func (b IgnoreBehavior) IntoExpr() ast.Expr {
	return ast.NewIdentifier(b.field)
}

type OverwriteBehavior struct {
	field string
}

func Overwrite(field string) OverwriteBehavior {
	return OverwriteBehavior{
		field: field,
	}
}

func (b OverwriteBehavior) Field() string { return b.field }

func (b OverwriteBehavior) IntoExpr() ast.Expr {
	return ast.NewValuesLiteral(
		ast.NewIdentifier(b.field),
	)
}
