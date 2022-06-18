package sql

import (
	"github.com/mpmaia/goneat/drivers/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlParser_GetColumnstColumns(t *testing.T) {

	parser := SqlParser{driver: sqlite.NewSqliteDriver()}
	columns, err := parser.GetColumns("SELECT a \"a1\", b, c \"c1\" FROM table")

	assert.NoError(t, err)
	assert.Equal(t, columns[0].Name, "a")
	assert.Equal(t, *columns[0].Alias, "a1")
	assert.Equal(t, columns[1].Name, "b")
	assert.Equal(t, columns[2].Name, "c")
	assert.Equal(t, *columns[2].Alias, "c1")

	columns, err = parser.GetColumns("UPDATE table SET a='a1', b=2, c=2.2")
	assert.NoError(t, err)
	assert.Equal(t, columns[0].Name, "a")
	assert.Equal(t, columns[1].Name, "b")
	assert.Equal(t, columns[2].Name, "c")
}
