package ast

type OnDuplicateKey struct {
	KeyIdents []*Identifier
	Updates   []*BinaryExpr
}

func NewOnDuplicateKey(keyParts ...*Identifier) *OnDuplicateKey {
	return &OnDuplicateKey{
		KeyIdents: keyParts,
	}
}

func (odk *OnDuplicateKey) Update(ident *Identifier, val IntoExpr) {
	odk.Updates = append(odk.Updates, NewBinaryExpr(ident, BinaryEquals, val))
}

func (odk *OnDuplicateKey) AcceptVisitor(fn func(n Node) bool) {
	if fn(odk) {
		for _, ident := range odk.KeyIdents {
			ident.AcceptVisitor(fn)
		}
		for _, u := range odk.Updates {
			u.AcceptVisitor(fn)
		}
	}
}

type ValuesLiteral struct {
	Expr
	Target *Identifier
}

func NewValuesLiteral(target *Identifier) *ValuesLiteral {
	return &ValuesLiteral{
		Target: target,
	}
}

func (vl *ValuesLiteral) AcceptVisitor(fn func(n Node) bool) {
	if fn(vl) {
		vl.Target.AcceptVisitor(fn)
	}
}
