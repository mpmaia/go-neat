package sql

import (
	"bytes"
	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/mpmaia/goneat/drivers"
	"strings"
)

type SqlParser struct {
	driver drivers.NeatDriver
}

type SqlColumn struct {
	Name  string
	Alias *string
}

func unescape(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}

func (p *SqlParser) GetColumns(sqlStmt string) ([]SqlColumn, error) {

	parser, err := xsqlparser.NewParser(bytes.NewBufferString(sqlStmt), p.driver.GetDialect())
	if err != nil {
		return nil, err
	}

	stmt, err := parser.ParseStatement()
	if err != nil {
		return nil, err
	}

	var columns []SqlColumn

	sqlast.Inspect(stmt, func(node sqlast.Node) bool {
		switch node.(type) {
		case nil:
			return false
		case *sqlast.AliasSelectItem:
			alias := unescape(node.(*sqlast.AliasSelectItem).Alias.Value)
			name := node.(*sqlast.AliasSelectItem).Expr.(*sqlast.Ident).Value
			columns = append(columns, SqlColumn{Alias: &alias, Name: name})
			return true
		case *sqlast.UnnamedSelectItem:
			name := node.(*sqlast.UnnamedSelectItem).Node.(*sqlast.Ident).Value
			columns = append(columns, SqlColumn{Name: name})
			return true
		case *sqlast.Assignment:
			name := node.(*sqlast.Assignment).ID.Value
			columns = append(columns, SqlColumn{Name: name})
			return true
		default:
			return true
		}
	})
	return columns, nil
}
