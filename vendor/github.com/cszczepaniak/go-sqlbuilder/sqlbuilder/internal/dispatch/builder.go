package dispatch

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"

type builder interface {
	Build() (statement.Statement, error)
}
