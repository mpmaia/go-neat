package go_neat

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

// Converts a StructField to a column object
func fieldToColumn(field reflect.StructField) (*column, error) {
	column := new(column)
	column.fieldName = field.Name
	column.dbName = strings.ToLower(field.Name)
	if tag, ok := field.Tag.Lookup("neat"); ok {
		parts := strings.Split(tag, ",")
		for i := 0; i < len(parts); i++ {
			if len(strings.TrimSpace(parts[i])) == 0 {
				continue
			}
			switch i {
			case 0:
				column.dbName = parts[i]
			case 1:
				column.decl = parts[i]
			}
		}
	} else {
		return nil, fmt.Errorf("Field %s does not have `neat` tag.", field.Name)
	}
	return column, nil
}

// Converts a reflect.Type to a slice of column objects
func typeToColumns(t reflect.Type) []*column {
	columns := make([]*column, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if c, err := fieldToColumn(t.Field(i)); err == nil {
			columns = append(columns, c)
		}
	}
	return columns
}

func modelToColumns(model interface{}) []*column {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		return modelToColumns(reflect.ValueOf(model).Elem().Interface())
	} else if typ.Kind() == reflect.Struct {
		return typeToColumns(typ)
	} else {
		panic("Unsupported type")
	}
}

// Finds the db column name that corresponds to the go struct's field name provided
func getDbNameByFieldName(column []*column, fieldName string) string {
	for _, c := range column {
		if c.fieldName == fieldName {
			return c.dbName
		}
	}
	return fieldName
}

// Gets the value from the field fieldName on the struct model. Expects a pointer to a struct or a struct.
func getValueFromColumn(model interface{}, fieldName string) reflect.Value {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		return getValueFromColumn(reflect.ValueOf(model).Elem().Interface(), fieldName)
	} else if typ.Kind() == reflect.Struct {
		return reflect.ValueOf(model).FieldByName(fieldName)
	} else {
		panic("Unsupported type")
	}
}

// Gets a pointer to the field called fieldName from the struct provided. Expects a pointer to a struct
func getPointerToColumn(dto interface{}, fieldName string) reflect.Value {
	return reflect.ValueOf(dto).Elem().FieldByName(fieldName).Addr()
}

// Converts a struct to sql named params
func dtoToSqlArgs(model interface{}) []interface{} {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Struct {
		flds := typeToColumns(typ)
		namedParams := make([]interface{}, len(flds))
		for i, fld := range flds {
			v := getValueFromColumn(model, fld.fieldName).Interface()
			namedParams[i] = sql.Named(fld.dbName, v)
		}
		return namedParams
	} else {
		return []interface{}{model}
	}
}

// Joins the column names using the char provided
func joinColumnNames(columns []*column, char byte) string {
	var sb strings.Builder
	for i, c := range columns {
		sb.WriteString(fmt.Sprintf("\"%s\"", c.dbName))
		if i < len(columns)-1 {
			sb.WriteByte(char)
		}
	}
	return sb.String()
}

func getValuesFromColumns(columns []*column, model interface{}) []interface{} {
	values := make([]interface{}, len(columns))
	for i, c := range columns {
		v := getValueFromColumn(model, c.fieldName)
		values[i] = v.Interface()
	}
	return values
}

func selectPointersToValuesFromColumns(columns []*column, cols []string, model interface{}) []interface{} {
	values := make([]interface{}, len(cols))
	for i, columnName := range cols {
		for _, column := range columns {
			if column.dbName == columnName {
				values[i] = reflect.ValueOf(model).Elem().FieldByName(column.fieldName).Addr().Interface()
				break
			}
		}
	}
	return values
}

func joinColumnsCustom(columns []*column, sb *strings.Builder, joinChar byte, cb func(int, *column) string) {
	for i, c := range columns {
		sb.WriteString(cb(i, c))
		if i < len(columns)-1 {
			sb.WriteByte(joinChar)
		}
	}
}

func getTempPath(fileName string) string {
	dir, err := ioutil.TempDir("", "neatdb-test")
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, fileName)
}
