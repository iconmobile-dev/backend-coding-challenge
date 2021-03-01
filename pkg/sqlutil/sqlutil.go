package sqlutil

import (
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/iconmobile-dev/go-core/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Select returns a postgres flavored SelectBuilder
func Select(columns ...string) sq.SelectBuilder {
	return psql.Select(columns...)
}

// GetColumns returns the columns names a structs fields are mapped to
func GetColumns(s interface{}) ([]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, errors.E(fmt.Errorf("failed to get columns as it is not a struct"), errors.Internal)
	}

	var columns []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !v.Field(i).CanInterface() {
			continue
		}

		column, ok := field.Tag.Lookup("db")
		if ok && column != "-" {
			columns = append(columns, column)
		}
		if !ok && column != "-" {
			column = strings.ToLower(field.Name)
			columns = append(columns, column)
		}
	}

	return columns, nil
}

// GetColumnMapping returns the a mapping from db columns names and struct field names to db columns names
func GetColumnMapping(s interface{}) (map[string]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, errors.E(fmt.Errorf("failed to get columns as it is not a struct"), errors.Internal)
	}

	m := map[string]string{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !v.Field(i).CanInterface() {
			continue
		}

		column, ok := field.Tag.Lookup("db")
		if ok && column == "-" {
			continue
		}
		if !ok {
			column = strings.ToLower(field.Name)
		}

		m[column] = column
		m[field.Name] = column
	}

	return m, nil
}

// GetColumnAliases returns the columns names a struct's fields are mapped to
func GetColumnAliases(columnPrefix string, s interface{}) ([]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, errors.E(fmt.Errorf("failed to get columns as it is not a struct"), errors.Internal)
	}

	var columns []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !v.Field(i).CanInterface() || field.Type.Kind() == reflect.Struct {
			continue
		}

		column, ok := field.Tag.Lookup("db")
		if ok && column != "-" {
			columns = append(columns, fmt.Sprintf(`%v%v "%v%v"`, columnPrefix, column, columnPrefix, column))
		}
		if !ok && column != "-" {
			column = strings.ToLower(field.Name)
			columns = append(columns, fmt.Sprintf(`%v%v "%v%v"`, columnPrefix, column, columnPrefix, column))
		}
	}

	return columns, nil
}

// ConcatSelectColumns returns the columns names a structs fields are mapped to
func ConcatSelectColumns(columns string, columnSlices ...[]string) string {
	var columnSlice []string
	if columns != "" {
		columnSlice = append(columnSlice, columns)
	}

	for _, cs := range columnSlices {
		columnSlice = append(columnSlice, cs...)
	}

	return strings.Join(columnSlice, ",\n")
}
