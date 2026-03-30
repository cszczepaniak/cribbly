package ast

type Insert struct {
	Into           TableExpr
	Columns        []*Identifier
	Values         []*TupleLiteral
	OnDuplicateKey *OnDuplicateKey
}

func NewInsert(into TableExpr, cols ...*Identifier) *Insert {
	return &Insert{
		Into:    into,
		Columns: cols,
	}
}

func (s *Insert) AddValues(vals ...IntoExpr) {
	s.Values = append(s.Values, NewTupleLiteral(vals...))
}

func (s *Insert) OnDuplicateKeyUpdate(keyParts []*Identifier, ident *Identifier, val IntoExpr) {
	if s.OnDuplicateKey == nil {
		s.OnDuplicateKey = NewOnDuplicateKey(keyParts...)
	}
	s.OnDuplicateKey.Update(ident, val)
}

func (s *Insert) AcceptVisitor(fn func(n Node) bool) {
	if fn(s) {
		s.Into.AcceptVisitor(fn)
		for _, col := range s.Columns {
			col.AcceptVisitor(fn)
		}
		for _, v := range s.Values {
			v.AcceptVisitor(fn)
		}
		if s.OnDuplicateKey != nil {
			s.OnDuplicateKey.AcceptVisitor(fn)
		}
	}
}
