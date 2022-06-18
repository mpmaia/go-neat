// Package goneat implements a simple ORM for Go inspired on the .NET Dapper library
package go_neat

import (
	"database/sql"
	"fmt"
	"github.com/mpmaia/goneat/drivers"
	reflection "github.com/mpmaia/goneat/pkg/reflection"
	"reflect"
	"strings"
)

// NeatDB is a handle that abstracts a database connection. See sql.DB for more information
type NeatDB struct {
	*sql.DB
	driver drivers.NeatDriver
}

// NeatRows contains the result of a query. See sql.Rows for more information
type NeatRows struct {
	*sql.Rows
}

// Open opens a connection to a database.
// This function follows the same API of the sql.DB.Open() function
func Open(driver drivers.NeatDriver, dsn string) (*NeatDB, error) {
	pool, err := sql.Open(driver.GetSqlDriverName(), dsn)
	return &NeatDB{pool, driver}, err
}

// Creates a table using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatCreateTable(model interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := reflection.TypeToColumns(typ)
	//constructs CREATE TABLE statement
	sql, err := db.driver.GetCreateTableStmt(tableName, columns)
	if err != nil {
		return nil, err
	}
	//creates the table
	return db.Exec(sql)
}

// Inserts a row mapping the columns using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatInsert(model interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := reflection.TypeToColumns(typ)
	values := reflection.GetValuesFromColumns(columns, model)
	//constructs INSERT statement
	sql, err := db.driver.GetInsertStmt(tableName, columns)
	if err != nil {
		return nil, err
	}
	//creates the table
	return db.Exec(sql, values...)
}

// Updates a row mapping the columns using the metadata on the struct field's tag provided on the parameter model
func (db *NeatDB) NeatUpdate(model interface{}, keyField string, keyValue interface{}) (sql.Result, error) {
	typ := reflect.TypeOf(model)
	tableName := strings.ToUpper(typ.Name())
	columns := reflection.TypeToColumns(typ)
	values := reflection.GetValuesFromColumns(columns, model)
	//constructs UPDATE statement
	sql, err := db.driver.GetUpdateStmt(tableName, columns, keyField)
	if err != nil {
		return nil, err
	}
	values = append(values, keyValue)
	//executes update
	return db.Exec(sql, values...)
}

// Selects a set of rows from the database mapping the columns of the result set to the struct returned by the factory parameter.
// This method requires that the struct has the neat tags declared with the expected db column name
func (db *NeatDB) NeatSelectAll(query string, factory func() interface{}, args ...interface{}) ([]interface{}, error) {
	model := factory()
	columns := reflection.ModelToColumns(model)
	if rows, err := db.Query(query, args...); err == nil {
		result := make([]interface{}, 0)
		if cols, err := rows.Columns(); err == nil {
			for rows.Next() {
				values := reflection.SelectPointersToValuesFromColumns(columns, cols, model)
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
