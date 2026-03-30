package ast

type CreateTable struct {
	Name        *Identifier
	IfNotExists bool

	Columns    []*ColumnSpec
	PrimaryKey *PrimaryKey
}

func NewCreateTable(name string) *CreateTable {
	return &CreateTable{
		Name: NewIdentifier(name),
	}
}

func (c *CreateTable) AcceptVisitor(fn func(n Node) bool) {
	if fn(c) {
		c.Name.AcceptVisitor(fn)
		for _, col := range c.Columns {
			col.AcceptVisitor(fn)
		}
	}
}

func (c *CreateTable) CreateIfNotExists() {
	c.IfNotExists = true
}

func (c *CreateTable) AddColumn(cs *ColumnSpec) {
	c.Columns = append(c.Columns, cs)
	if cs.ComprisesPrimaryKey {
		c.addPrimaryKeyColumn(cs.Name.Name)
	}
}

func (c *CreateTable) addPrimaryKeyColumn(colName string) {
	if c.PrimaryKey == nil {
		c.PrimaryKey = NewPrimaryKey()
	}

	c.PrimaryKey.AddColumn(colName)
}
