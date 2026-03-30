package column

import (
	"time"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type tinyIntColumnBuilder struct {
	*integerColumnBuilder[int8, *tinyIntColumnBuilder]
}

func (*tinyIntColumnBuilder) columnType() ast.ColumnType {
	return ast.TinyInt()
}

func TinyInt(name string) *tinyIntColumnBuilder {
	b := &tinyIntColumnBuilder{}
	b.integerColumnBuilder = newIntegerColumnBuilder[int8](name, b)
	return b
}

type smallIntColumnBuilder struct {
	*integerColumnBuilder[int16, *smallIntColumnBuilder]
}

func (*smallIntColumnBuilder) columnType() ast.ColumnType {
	return ast.SmallInt()
}

func SmallInt(name string) *smallIntColumnBuilder {
	b := &smallIntColumnBuilder{}
	b.integerColumnBuilder = newIntegerColumnBuilder[int16](name, b)
	return b
}

type intColumnBuilder struct {
	*integerColumnBuilder[int32, *intColumnBuilder]
}

func (*intColumnBuilder) columnType() ast.ColumnType {
	return ast.Int()
}

func Int(name string) *intColumnBuilder {
	b := &intColumnBuilder{}
	b.integerColumnBuilder = newIntegerColumnBuilder[int32](name, b)
	return b
}

type bigIntColumnBuilder struct {
	*integerColumnBuilder[int64, *bigIntColumnBuilder]
}

func (*bigIntColumnBuilder) columnType() ast.ColumnType {
	return ast.BigInt()
}

func BigInt(name string) *bigIntColumnBuilder {
	b := &bigIntColumnBuilder{}
	b.integerColumnBuilder = newIntegerColumnBuilder[int64](name, b)
	return b
}

type charColumnBuilder struct {
	*stringColumnBuilder[*charColumnBuilder]
}

func (b *charColumnBuilder) columnType() ast.ColumnType {
	return ast.Char(b.size)
}

func Char(name string, size int) *charColumnBuilder {
	b := &charColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type varCharColumnBuilder struct {
	*stringColumnBuilder[*varCharColumnBuilder]
}

func (b *varCharColumnBuilder) columnType() ast.ColumnType {
	return ast.VarChar(b.size)
}

func VarChar(name string, size int) *varCharColumnBuilder {
	b := &varCharColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type textColumnBuilder struct {
	*stringColumnBuilder[*textColumnBuilder]
}

func (b *textColumnBuilder) columnType() ast.ColumnType {
	return ast.Text(b.size)
}

func Text(name string, size int) *textColumnBuilder {
	b := &textColumnBuilder{}
	b.stringColumnBuilder = newStringColumnBuilder(name, size, b)
	return b
}

type tinyBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *tinyBlobColumnBuilder]
}

func (b *tinyBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.TinyBlob()
}

func TinyBlob(name string) *tinyBlobColumnBuilder {
	b := &tinyBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type blobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *blobColumnBuilder]
}

func (b *blobColumnBuilder) columnType() ast.ColumnType {
	return ast.Blob()
}

func Blob(name string) *blobColumnBuilder {
	b := &blobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type mediumBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *mediumBlobColumnBuilder]
}

func (b *mediumBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.MediumBlob()
}

func MediumBlob(name string) *mediumBlobColumnBuilder {
	b := &mediumBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type longBlobColumnBuilder struct {
	*baseColumnBuilder[[]byte, *longBlobColumnBuilder]
}

func (b *longBlobColumnBuilder) columnType() ast.ColumnType {
	return ast.LongBlob()
}

func LongBlob(name string) *longBlobColumnBuilder {
	b := &longBlobColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[[]byte](name, b)
	return b
}

type dateTimeColumnBuilder struct {
	*baseColumnBuilder[time.Time, *dateTimeColumnBuilder]
}

func (b *dateTimeColumnBuilder) columnType() ast.ColumnType {
	return ast.DateTime()
}

func DateTime(name string) *dateTimeColumnBuilder {
	b := &dateTimeColumnBuilder{}
	b.baseColumnBuilder = newBaseColumnBuilder[time.Time](name, b)
	return b
}
