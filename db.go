// Package goneat implements a simple ORM for Go inspired on the .NET Dapper library
package go_neat

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Describes the column metadata
type column struct {
	fieldName string
	dbName    string
	decl      string
}

// A handle describing a database connection. See sql.DB for more information
type NeatDB struct {
	*sql.DB
}

// Contains the result of a query. See sql.Rows for more information
type NeatRows struct {
	*sql.Rows
}

//prints the content of the column
func (c *column) String() string {
	return fmt.Sprintf("fieldName:%s,dbName:%s,decl:%s", c.fieldName, c.dbName, c.decl)
}

// Opens a connection to a database.
// This function follows the same API of the sql.DB.Open() function
func Open(driverName string, dsn string) (*NeatDB, error) {
	pool, err := sql.Open(driverName, dsn)
	return &NeatDB{pool}, err
}

// Creates a table using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatCreateTable(model interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := typeToColumns(typ)
	//constructs CREATE TABLE statement
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))
	for i, c := range columns {
		sb.WriteString(fmt.Sprintf("\"%s\" %s", c.dbName, c.decl))
		if i < len(columns)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(");\n")
	//creates the table
	return db.Exec(sb.String())
}

// Inserts a row mapping the columns using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatInsert(model interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := typeToColumns(typ)
	values := getValuesFromColumns(columns, model)
	//constructs CREATE TABLE statement
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO \"%s\"\n (", tableName))
	sb.WriteString(joinColumnNames(columns, ','))
	sb.WriteString(") VALUES (")
	joinColumnsCustom(columns, &sb, ',', func(i int, c *column) string {
		return "?"
	})
	sb.WriteString(");")
	//creates the table
	return db.Exec(sb.String(), values...)
}

// Updates a row mapping the columns using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatUpdate(model interface{}, keyField string, keyValue interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := typeToColumns(typ)
	values := getValuesFromColumns(columns, model)
	//constructs UPDATE statement
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UPDATE \"%s\" SET \n", tableName))
	for i, c := range columns {
		sb.WriteString(c.dbName + " = ?")
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(fmt.Sprintf(" WHERE %s = ?", getDbNameByFieldName(columns, keyField)))
	values = append(values, keyValue)
	//executes update
	return db.Exec(sb.String(), values...)
}

// Selects a set of rows from the database mapping the columns of the result set to the struct returned by the factory parameter.
// This method requires that the struct has the neat tags declared with the expected db column name
func (db *NeatDB) NeatSelectAll(query string, factory func() interface{}, args ...interface{}) ([]interface{}, error) {
	model := factory()
	columns := modelToColumns(model)
	if rows, err := db.Query(query, args...); err == nil {
		result := make([]interface{}, 0)
		if cols, err := rows.Columns(); err == nil {
			for rows.Next() {
				values := selectPointersToValuesFromColumns(columns, cols, model)
				rows.Scan(values...)
				result = append(result, model)
				model = factory()
			}
		} else {
			return nil, err
		}
		return result, nil
	} else {
		return nil, err
	}
}

// Selects a unique row from the database mapping the columns of the result to the struct returned by the factory parameter.
// This method requires that the struct has the neat tags declared with the expected db column name
func (db *NeatDB) NeatSelectOne(query string, factory func() interface{}, args ...interface{}) (interface{}, error) {
	if result, err := db.NeatSelectAll(query, factory, args...); err == nil {
		if len(result) > 0 {
			return result[0], nil
		} else {
			return nil, fmt.Errorf("Record not found for query %s and args %s", query, args)
		}
	} else {
		return nil, err
	}
}
