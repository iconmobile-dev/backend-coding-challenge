package sqlutil

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iconmobile-dev/go-core/errors"

	sq "github.com/Masterminds/squirrel"
)

// TimeFilter specifies filter criteria for an timestamp column
type TimeFilter struct {
	Before *time.Time
	After  *time.Time
}

// UseTimeFilter adds filter criteria defined in IntFilter to a column in a sql query
func UseTimeFilter(q sq.SelectBuilder, column string, filter TimeFilter) sq.SelectBuilder {
	if filter.Before != nil {
		q = q.Where(sq.Lt{column: filter.Before})
	}

	if filter.After != nil {
		q = q.Where(sq.Gt{column: filter.After})
	}

	return q
}

// IntFilter specifies filter criteria for an integer column
type IntFilter struct {
	Is    *int
	Not   *int
	In    []int
	NotIn []int
	Gt    *int
	Gte   *int
	Lt    *int
	Lte   *int
}

// UseIntFilter adds filter criteria defined in IntFilter to a column in a sql query
func UseIntFilter(q sq.SelectBuilder, column string, filter IntFilter) sq.SelectBuilder {
	if filter.Is != nil {
		q = q.Where(sq.Eq{column: filter.Is})
	}

	if filter.Not != nil {
		q = q.Where(sq.NotEq{column: filter.Not})
	}

	if filter.In != nil {
		q = q.Where(sq.Eq{column: filter.In})
	}

	if filter.NotIn != nil {
		q = q.Where(sq.NotEq{column: filter.NotIn})
	}

	if filter.Gt != nil {
		q = q.Where(sq.Gt{column: filter.Gt})
	}

	if filter.Gte != nil {
		q = q.Where(sq.GtOrEq{column: filter.Gte})
	}

	if filter.Lt != nil {
		q = q.Where(sq.Lt{column: filter.Lt})
	}

	if filter.Lte != nil {
		q = q.Where(sq.LtOrEq{column: filter.Lte})
	}

	return q
}

// StringFilter specifies filter criteria for a string column
type StringFilter struct {
	CaseSensitive *bool
	Is            *string
	Not           *string
	In            []string
	NotIn         []string
	Contains      *string
	NotContains   *string
	StartsWith    *string
	NotStartsWith *string
	EndsWith      *string
	NotEndsWith   *string
}

// UseStringFilter adds filter criteria defined in StringFilter to a column in a sql query
func UseStringFilter(q sq.SelectBuilder, column string, filter StringFilter) sq.SelectBuilder {
	if filter.Is != nil {
		q = q.Where(sq.Eq{column: filter.Is})
	}

	if filter.Not != nil {
		q = q.Where(sq.NotEq{column: filter.Not})
	}

	if filter.In != nil {
		q = q.Where(sq.Eq{column: filter.In})
	}

	if filter.NotIn != nil {
		q = q.Where(sq.NotEq{column: filter.NotIn})
	}

	if filter.Contains != nil {
		q = q.Where(sq.ILike{column: "%" + *filter.Contains + "%"})
	}

	if filter.NotContains != nil {
		q = q.Where(sq.NotILike{column: "%" + *filter.NotContains + "%"})
	}

	if filter.StartsWith != nil {
		q = q.Where(sq.Like{column: *filter.StartsWith + "%"})
	}

	if filter.NotStartsWith != nil {
		q = q.Where(sq.NotLike{column: *filter.NotStartsWith + "%"})
	}

	if filter.EndsWith != nil {
		q = q.Where(sq.Like{column: "%" + *filter.EndsWith})
	}

	if filter.NotEndsWith != nil {
		q = q.Where(sq.NotLike{column: "%" + *filter.NotEndsWith})
	}

	return q
}

// BoolFilter specifies filter criteria for a boolean column
type BoolFilter struct {
	Is *bool
}

// UseBoolFilter adds filter criteria defined in BoolFilter to a column in a sql query
func UseBoolFilter(q sq.SelectBuilder, column string, filter BoolFilter) sq.SelectBuilder {
	if filter.Is != nil {
		q = q.Where(sq.Eq{column: filter.Is})
	}

	return q
}

// UseStructFilter adds filter criteria defined in a struct filter to columns in a sql query
// A struct filter may look like the following:
//
//	type testRowFilter struct {
//		ID           *sqlutil.IntFilter
//		TestCase     *sqlutil.StringFilter `db:"test_case"`
//		IntColumn    *sqlutil.IntFilter    `db:"int_column"`
//		StringColumn *sqlutil.StringFilter `db:"string_column"`
//		BoolColumn   *sqlutil.BoolFilter   `db:"bool_column"`
//		TimeColumn   *sqlutil.TimeFilter   `db:"time_column"`
//		unexported   interface{}
//	}
//
func UseStructFilter(q sq.SelectBuilder, columnPrefix string, filter interface{}) (sq.SelectBuilder, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() != reflect.Struct {
		return q, errors.E(fmt.Errorf("failed to use filter as it is not a struct"), errors.Internal)
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !v.Field(i).CanInterface() {
			continue
		}

		column, ok := field.Tag.Lookup("db")
		if !ok {
			column = strings.ToLower(field.Name)
		}
		column = columnPrefix + column

		switch f := v.Field(i).Interface().(type) {
		case *BoolFilter:
			if f != nil {
				q = UseBoolFilter(q, column, *f)
			}
		case *StringFilter:
			if f != nil {
				q = UseStringFilter(q, column, *f)
			}
		case *TimeFilter:
			if f != nil {
				q = UseTimeFilter(q, column, *f)
			}
		case *IntFilter:
			if f != nil {
				q = UseIntFilter(q, column, *f)
			}
		}
	}

	return q, nil
}
