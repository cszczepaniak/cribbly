package ast

type ColumnType interface {
	Expr
	columnType()
}

type (
	TinyIntColumn  struct{ ColumnType }
	SmallIntColumn struct{ ColumnType }
	IntColumn      struct{ ColumnType }
	BigIntColumn   struct{ ColumnType }
	CharColumn     struct {
		ColumnType
		Size int
	}
	VarCharColumn struct {
		ColumnType
		Size int
	}
	TextColumn struct {
		ColumnType
		Size int
	}
	TinyBlobColumn   struct{ ColumnType }
	BlobColumn       struct{ ColumnType }
	MediumBlobColumn struct{ ColumnType }
	LongBlobColumn   struct{ ColumnType }
	DateTimeColumn   struct{ ColumnType }
)

func TinyInt() TinyIntColumn {
	return TinyIntColumn{}
}

func (c TinyIntColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func SmallInt() SmallIntColumn {
	return SmallIntColumn{}
}

func (c SmallIntColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func Int() IntColumn {
	return IntColumn{}
}

func (c IntColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func BigInt() BigIntColumn {
	return BigIntColumn{}
}

func (c BigIntColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func Char(size int) CharColumn {
	return CharColumn{
		Size: size,
	}
}

func (c CharColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func VarChar(size int) VarCharColumn {
	return VarCharColumn{
		Size: size,
	}
}

func (c VarCharColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func Text(size int) TextColumn {
	return TextColumn{
		Size: size,
	}
}

func (c TextColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func TinyBlob() TinyBlobColumn {
	return TinyBlobColumn{}
}

func (c TinyBlobColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func Blob() BlobColumn {
	return BlobColumn{}
}

func (c BlobColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func MediumBlob() MediumBlobColumn {
	return MediumBlobColumn{}
}

func (c MediumBlobColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func LongBlob() LongBlobColumn {
	return LongBlobColumn{}
}

func (c LongBlobColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}

func DateTime() DateTimeColumn {
	return DateTimeColumn{}
}

func (c DateTimeColumn) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}
