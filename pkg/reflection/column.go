package reflection

import "fmt"

// DBColumn holds column metadata
type DBColumn struct {
	FieldName string
	DbName    string
	Decl      string
}

type DBColumnList []*DBColumn

// prints the content of the column
func (c *DBColumn) String() string {
	return fmt.Sprintf("fieldName:%s,dbName:%s,decl:%s", c.FieldName, c.DbName, c.Decl)
}
