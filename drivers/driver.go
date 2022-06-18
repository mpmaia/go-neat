package drivers

import (
	"github.com/akito0107/xsqlparser/dialect"
	"github.com/mpmaia/goneat/pkg/reflection"
)

// NeatDriver abstracts database specific code
type NeatDriver interface {
	// GetDialect returns the sql dialect implementation for this database
	GetDialect() dialect.Dialect
	// GetSqlDriverName returns the sql driver name of this database
	GetSqlDriverName() string
	// GetCreateTableStmt constructs a creates table statement for the current database
	GetCreateTableStmt(tableName string, columns reflection.DBColumnList) (string, error)
	// GetUpdateStmt returns an update statement for the selected columns
	GetUpdateStmt(tableName string, columns reflection.DBColumnList, keyField string) (string, error)
	// GetInsertStmt returns an insert statement for the given columns
	GetInsertStmt(tableName string, columns reflection.DBColumnList) (string, error)
}
