package sqlutil_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/ptrutil"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/sqlutil"
	"github.com/stretchr/testify/assert"
)

func TestUseStringFilter(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	t.Run(`.Is == ".Is"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".Is")},
			{StringColumn: ptrutil.String(".Is")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			Is: ptrutil.String(".Is"),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.Not == ".Not"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".Not")},
			{StringColumn: ptrutil.String(".Not")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			Not: ptrutil.String(".Not"),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.In == [".InA", ".InB"]`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".InA")},
			{StringColumn: ptrutil.String(".InA")},
			{StringColumn: ptrutil.String(".InB")},
			{StringColumn: ptrutil.String(".InB")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
			{StringColumn: ptrutil.String("d")},
			{StringColumn: ptrutil.String("e")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			In: []string{".InA", ".InB"},
		})

		assert.Len(t, filtered, 4)
	})

	t.Run(`.NotIn == [".NotInA", ".NotInB"]`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".NotInA")},
			{StringColumn: ptrutil.String(".NotInA")},
			{StringColumn: ptrutil.String(".NotInB")},
			{StringColumn: ptrutil.String(".NotInB")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
			{StringColumn: ptrutil.String("d")},
			{StringColumn: ptrutil.String("e")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			NotIn: []string{".NotInA", ".NotInB"},
		})

		assert.Len(t, filtered, 5)
	})

	t.Run(`.Contains == ".Contains"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("__.Contains__")},
			{StringColumn: ptrutil.String(".Contains__")},
			{StringColumn: ptrutil.String("__.Contains")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
			{StringColumn: ptrutil.String("d")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			Contains: ptrutil.String(".Contains"),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.Contains == ".contains"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("__.Contains__")},
			{StringColumn: ptrutil.String(".Contains__")},
			{StringColumn: ptrutil.String("__.Contains")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
			{StringColumn: ptrutil.String("d")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			Contains: ptrutil.String(".contains"),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.NotContains == ".NotContains"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("__.NotContains__")},
			{StringColumn: ptrutil.String(".NotContains__")},
			{StringColumn: ptrutil.String("__.NotContains")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
			{StringColumn: ptrutil.String("d")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			NotContains: ptrutil.String(".NotContains"),
		})

		assert.Len(t, filtered, 4)
	})

	t.Run(`.StartsWith == ".StartsWith"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".StartsWithA")},
			{StringColumn: ptrutil.String(".StartsWithB")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			StartsWith: ptrutil.String(".StartsWith"),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.NotStartsWith == ".NotStartsWith"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".NotStartsWithA")},
			{StringColumn: ptrutil.String(".NotStartsWithB")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			NotStartsWith: ptrutil.String(".NotStartsWith"),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.EndsWith == ".EndsWith"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("A.EndsWith")},
			{StringColumn: ptrutil.String("B.EndsWith")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			EndsWith: ptrutil.String(".EndsWith"),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.NotEndsWith == ".NotEndsWith"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("A.NotEndsWith")},
			{StringColumn: ptrutil.String("B.NotEndsWith")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered := mustUseStringFilter(t, db, t.Name(), sqlutil.StringFilter{
			NotEndsWith: ptrutil.String(".NotEndsWith"),
		})

		assert.Len(t, filtered, 3)
	})
}

func mustUseStringFilter(t *testing.T, db *storage.DB, testCase string, filter sqlutil.StringFilter) []testRow {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q = sqlutil.UseStringFilter(q, "string_column", filter)

	sql, args, err := q.ToSql()
	require.NoError(t, err)

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	require.NoError(t, err)

	return rs
}

func TestUseIntFilter(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	t.Run(`.Is == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
			{IntColumn: ptrutil.Int(3)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Is: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.Not == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
			{IntColumn: ptrutil.Int(3)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Not: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.In == [1371, 1372]`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(1371)},
			{IntColumn: ptrutil.Int(1371)},
			{IntColumn: ptrutil.Int(1372)},
			{IntColumn: ptrutil.Int(1372)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
			{IntColumn: ptrutil.Int(3)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			In: []int{1371, 1372},
		})

		assert.Len(t, filtered, 4)
	})

	t.Run(`.NotIn == [1371, 1372]`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(1371)},
			{IntColumn: ptrutil.Int(1371)},
			{IntColumn: ptrutil.Int(1372)},
			{IntColumn: ptrutil.Int(1372)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
			{IntColumn: ptrutil.Int(3)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			NotIn: []int{1371, 1372},
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.Gt == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(134)},
			{IntColumn: ptrutil.Int(135)},
			{IntColumn: ptrutil.Int(136)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(138)},
			{IntColumn: ptrutil.Int(139)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Gt: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.Gte == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(134)},
			{IntColumn: ptrutil.Int(135)},
			{IntColumn: ptrutil.Int(136)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(138)},
			{IntColumn: ptrutil.Int(139)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Gte: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.Lt == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(134)},
			{IntColumn: ptrutil.Int(135)},
			{IntColumn: ptrutil.Int(136)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(138)},
			{IntColumn: ptrutil.Int(139)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Lt: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 3)
	})

	t.Run(`.Lte == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(134)},
			{IntColumn: ptrutil.Int(135)},
			{IntColumn: ptrutil.Int(136)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(138)},
			{IntColumn: ptrutil.Int(139)},
		})

		filtered := mustUseIntFilter(t, db, t.Name(), sqlutil.IntFilter{
			Lte: ptrutil.Int(137),
		})

		assert.Len(t, filtered, 4)
	})
}

func mustUseIntFilter(t *testing.T, db *storage.DB, testCase string, filter sqlutil.IntFilter) []testRow {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q = sqlutil.UseIntFilter(q, "int_column", filter)

	sql, args, err := q.ToSql()
	require.NoError(t, err)

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	require.NoError(t, err)

	return rs
}

func TestUseTimeFilter(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	t.Run(`.Before == time.Now()`, func(t *testing.T) {
		now := time.Now()

		mustInsertTestRows(t, db, t.Name(), []testRow{
			{TimeColumn: ptrutil.Time(now.Add(-2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(-1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now)},
			{TimeColumn: ptrutil.Time(now.Add(1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(3 * time.Millisecond))},
		})

		filtered := mustUseTimeFilter(t, db, t.Name(), sqlutil.TimeFilter{
			Before: ptrutil.Time(now),
		})

		assert.Len(t, filtered, 2)
	})

	t.Run(`.After == time.Now()`, func(t *testing.T) {
		now := time.Now()

		mustInsertTestRows(t, db, t.Name(), []testRow{
			{TimeColumn: ptrutil.Time(now.Add(-2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(-1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now)},
			{TimeColumn: ptrutil.Time(now.Add(1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(3 * time.Millisecond))},
		})

		filtered := mustUseTimeFilter(t, db, t.Name(), sqlutil.TimeFilter{
			After: ptrutil.Time(now),
		})

		assert.Len(t, filtered, 3)
	})
}

func mustUseTimeFilter(t *testing.T, db *storage.DB, testCase string, filter sqlutil.TimeFilter) []testRow {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q = sqlutil.UseTimeFilter(q, "time_column", filter)

	sql, args, err := q.ToSql()
	require.NoError(t, err)

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	require.NoError(t, err)

	return rs
}

func TestUseBoolFilter(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	t.Run(`.Is == true`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{BoolColumn: ptrutil.Bool(true)},
			{BoolColumn: ptrutil.Bool(true)},
			{BoolColumn: ptrutil.Bool(false)},
		})

		filtered := mustUseBoolFilter(t, db, t.Name(), sqlutil.BoolFilter{
			Is: ptrutil.Bool(true),
		})

		assert.Len(t, filtered, 2)
	})
}

func mustUseBoolFilter(t *testing.T, db *storage.DB, testCase string, filter sqlutil.BoolFilter) []testRow {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q = sqlutil.UseBoolFilter(q, "bool_column", filter)

	sql, args, err := q.ToSql()
	require.NoError(t, err)

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	require.NoError(t, err)

	return rs
}

type testRowFilter struct {
	ID           *sqlutil.IntFilter
	TestCase     *sqlutil.StringFilter `db:"test_case"`
	IntColumn    *sqlutil.IntFilter    `db:"int_column"`
	StringColumn *sqlutil.StringFilter `db:"string_column"`
	BoolColumn   *sqlutil.BoolFilter   `db:"bool_column"`
	TimeColumn   *sqlutil.TimeFilter   `db:"time_column"`
	unexported   interface{}
}

func TestUseStructFilter(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	t.Run(`.Id.Not == 0`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".Is")},
			{StringColumn: ptrutil.String(".Is")},
		})

		filtered, err := useStructFilter(t, db, t.Name(), testRowFilter{
			ID: &sqlutil.IntFilter{
				Not: ptrutil.Int(0),
			},
		})
		require.NoError(t, err)

		assert.Len(t, filtered, 2)
	})

	t.Run(`.StringColumn.Is == ".Is"`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String(".Is")},
			{StringColumn: ptrutil.String(".Is")},
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		filtered, err := useStructFilter(t, db, t.Name(), testRowFilter{
			StringColumn: &sqlutil.StringFilter{
				Is: ptrutil.String(".Is"),
			},
		})
		require.NoError(t, err)

		assert.Len(t, filtered, 2)
	})

	t.Run(`.BoolColumn.Is == true`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{BoolColumn: ptrutil.Bool(true)},
			{BoolColumn: ptrutil.Bool(true)},
			{BoolColumn: ptrutil.Bool(false)},
		})

		filtered, err := useStructFilter(t, db, t.Name(), testRowFilter{
			BoolColumn: &sqlutil.BoolFilter{
				Is: ptrutil.Bool(true),
			},
		})
		require.NoError(t, err)

		assert.Len(t, filtered, 2)
	})

	t.Run(`.TimeColumn.Before == time.Now()`, func(t *testing.T) {
		now := time.Now()

		mustInsertTestRows(t, db, t.Name(), []testRow{
			{TimeColumn: ptrutil.Time(now.Add(-2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(-1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now)},
			{TimeColumn: ptrutil.Time(now.Add(1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(2 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(3 * time.Millisecond))},
		})

		filtered, err := useStructFilter(t, db, t.Name(), testRowFilter{
			TimeColumn: &sqlutil.TimeFilter{
				Before: ptrutil.Time(now),
			},
		})
		require.NoError(t, err)

		assert.Len(t, filtered, 2)
	})

	t.Run(`.IntColumn.Is == 137`, func(t *testing.T) {
		mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(137)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
			{IntColumn: ptrutil.Int(3)},
		})

		filtered, err := useStructFilter(t, db, t.Name(), testRowFilter{
			IntColumn: &sqlutil.IntFilter{
				Is: ptrutil.Int(137),
			},
		})
		require.NoError(t, err)

		assert.Len(t, filtered, 2)
	})

	t.Run(`not a struct`, func(t *testing.T) {
		_, err := useStructFilter(t, db, t.Name(), "string")
		require.Error(t, err)
	})
}

func useStructFilter(t *testing.T, db *storage.DB, testCase string, filter interface{}) ([]testRow, error) {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q, err := sqlutil.UseStructFilter(q, "", filter)
	if err != nil {
		return nil, err
	}

	sql, args, err := q.ToSql()
	require.NoError(t, err)

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	require.NoError(t, err)

	return rs, nil
}
