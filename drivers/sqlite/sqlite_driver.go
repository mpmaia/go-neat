package sqlite

import (
	"fmt"
	"github.com/akito0107/xsqlparser/dialect"
	"github.com/mpmaia/goneat/drivers"
	"github.com/mpmaia/goneat/pkg/reflection"
	"strings"
)

type SqliteDriver struct {
}

func NewSqliteDriver() drivers.NeatDriver {
	return &SqliteDriver{}
}

func (d *SqliteDriver) GetDialect() dialect.Dialect {
	return &dialect.GenericSQLDialect{}
}

func (d *SqliteDriver) GetSqlDriverName() string {
	return "sqlite"
}

func (d *SqliteDriver) GetCreateTableStmt(tableName string, columns reflection.DBColumnList) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))
	for i, c := range columns {
		sb.WriteString(fmt.Sprintf("\"%s\" %s", c.DbName, c.Decl))
		if i < len(columns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(");\n")
	return sb.String(), nil
}

func (d *SqliteDriver) GetUpdateStmt(tableName string, columns reflection.DBColumnList, keyField string) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE \"%s\" SET \n", tableName))
	for i, c := range columns {
		sb.WriteString(c.DbName + " = ?")
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(fmt.Sprintf(" WHERE %s = ?", reflection.GetDbNameByFieldName(columns, keyField)))
	return sb.String(), nil
}

func (d *SqliteDriver) GetInsertStmt(tableName string, columns reflection.DBColumnList) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO \"%s\"\n (", tableName))
	sb.WriteString(reflection.JoinColumnNames(columns, ','))
	sb.WriteString(") VALUES (")
	reflection.JoinColumnsCustom(columns, &sb, ',', func(i int, c *reflection.DBColumn) string {
		return "?"
	})
	sb.WriteString(");")
	return sb.String(), nil
}
