package reflection

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// FieldToColumn converts a StructField to a DBColumn object
func FieldToColumn(field reflect.StructField) (*DBColumn, error) {
	column := new(DBColumn)
	column.FieldName = field.Name
	column.DbName = strings.ToLower(field.Name)
	if tag, ok := field.Tag.Lookup("neat"); ok {
		parts := strings.Split(tag, ",")
		for i := 0; i < len(parts); i++ {
			if len(strings.TrimSpace(parts[i])) == 0 {
				continue
			}
			switch i {
			case 0:
				column.DbName = parts[i]
			case 1:
				column.Decl = parts[i]
			}
		}
	} else {
		return nil, fmt.Errorf("Field %s does not have `neat` tag.", field.Name)
	}
	return column, nil
}

// TypeToColumns converts a reflect.Type to a slice of DBColumn objects
func TypeToColumns(t reflect.Type) []*DBColumn {
	columns := make([]*DBColumn, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if c, err := FieldToColumn(t.Field(i)); err == nil {
			columns = append(columns, c)
		}
	}
	return columns
}

// ModelToColumns extracts the DBColumn information from the model
func ModelToColumns(model interface{}) []*DBColumn {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		return ModelToColumns(reflect.ValueOf(model).Elem().Interface())
	} else if typ.Kind() == reflect.Struct {
		return TypeToColumns(typ)
	} else {
		panic("Unsupported type")
	}
}

// GetDbNameByFieldName finds the db DBColumn name that corresponds to the go struct's field name provided
func GetDbNameByFieldName(column []*DBColumn, fieldName string) string {
	for _, c := range column {
		if c.FieldName == fieldName {
			return c.DbName
		}
	}
	return fieldName
}

// GetValueFromColumn gets the value from the field FieldName on the struct model. Expects a pointer to a struct or a struct.
func GetValueFromColumn(model interface{}, fieldName string) reflect.Value {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		return GetValueFromColumn(reflect.ValueOf(model).Elem().Interface(), fieldName)
	} else if typ.Kind() == reflect.Struct {
		return reflect.ValueOf(model).FieldByName(fieldName)
	} else {
		panic("Unsupported type")
	}
}

// GetPointerToColumn gets a pointer to the field called FieldName from the struct provided. Expects a pointer to a struct
func GetPointerToColumn(dto interface{}, fieldName string) reflect.Value {
	return reflect.ValueOf(dto).Elem().FieldByName(fieldName).Addr()
}

// DtoToSqlArgs converts a struct to sql named params
func DtoToSqlArgs(model interface{}) []interface{} {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Struct {
		flds := TypeToColumns(typ)
		namedParams := make([]interface{}, len(flds))
		for i, fld := range flds {
			v := GetValueFromColumn(model, fld.FieldName).Interface()
			namedParams[i] = sql.Named(fld.DbName, v)
		}
		return namedParams
	} else {
		return []interface{}{model}
	}
}

// JoinColumnNames joins the DBColumn names using the char provided
func JoinColumnNames(columns []*DBColumn, char byte) string {
	var sb strings.Builder
	for i, c := range columns {
		sb.WriteString(fmt.Sprintf("\"%s\"", c.DbName))
		if i < len(columns)-1 {
			sb.WriteByte(char)
		}
	}
	return sb.String()
}

// GetValuesFromColumns extracts the values from each column from the model provided
func GetValuesFromColumns(columns []*DBColumn, model interface{}) []interface{} {
	values := make([]interface{}, len(columns))
	for i, c := range columns {
		v := GetValueFromColumn(model, c.FieldName)
		values[i] = v.Interface()
	}
	return values
}

// SelectPointersToValuesFromColumns gets pointers for each column value
func SelectPointersToValuesFromColumns(columns []*DBColumn, cols []string, model interface{}) []interface{} {
	values := make([]interface{}, len(cols))
	for i, columnName := range cols {
		for _, column := range columns {
			if column.DbName == columnName {
				values[i] = reflect.ValueOf(model).Elem().FieldByName(column.FieldName).Addr().Interface()
				break
			}
		}
	}
	return values
}

// JoinColumnsCustom Join column names
func JoinColumnsCustom(columns []*DBColumn, sb *strings.Builder, joinChar byte, cb func(int, *DBColumn) string) {
	for i, c := range columns {
		sb.WriteString(cb(i, c))
		if i < len(columns)-1 {
			sb.WriteByte(joinChar)
		}
	}
}
