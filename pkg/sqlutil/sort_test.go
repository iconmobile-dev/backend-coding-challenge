package sqlutil_test

import (
	"testing"
	"time"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/ptrutil"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/sqlutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseOneColumnSort(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	columnMapping, err := sqlutil.GetColumnMapping(testRow{})
	require.NoError(t, err)

	t.Run(`.Column == "" and .Order == ""`, func(t *testing.T) {
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(0)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[0].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[2].ID, sorted[2].ID)
	})

	t.Run(`.Column == "int_column" and .Order == "desc"`, func(t *testing.T) {
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(0)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "int_column",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`.Column == "IntColumn" and .Order == "desc"`, func(t *testing.T) {
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{IntColumn: ptrutil.Int(0)},
			{IntColumn: ptrutil.Int(1)},
			{IntColumn: ptrutil.Int(2)},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "IntColumn",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`.Column == "string_column" and .Order == "desc"`, func(t *testing.T) {
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "string_column",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`.Column == "StringColumn" and .Order == "desc"`, func(t *testing.T) {
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{StringColumn: ptrutil.String("a")},
			{StringColumn: ptrutil.String("b")},
			{StringColumn: ptrutil.String("c")},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "StringColumn",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`.Column == "time_column" and .Order == "desc"`, func(t *testing.T) {
		now := time.Now()
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{TimeColumn: ptrutil.Time(now)},
			{TimeColumn: ptrutil.Time(now.Add(1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(2 * time.Millisecond))},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "time_column",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`.Column == "TimeColumn" and .Order == "desc"`, func(t *testing.T) {
		now := time.Now()
		inserted := mustInsertTestRows(t, db, t.Name(), []testRow{
			{TimeColumn: ptrutil.Time(now)},
			{TimeColumn: ptrutil.Time(now.Add(1 * time.Millisecond))},
			{TimeColumn: ptrutil.Time(now.Add(2 * time.Millisecond))},
		})

		sorted, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{
			Column: "TimeColumn",
			Order:  "desc",
		}, columnMapping)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, inserted[2].ID, sorted[0].ID)
		assert.Equal(t, inserted[1].ID, sorted[1].ID)
		assert.Equal(t, inserted[0].ID, sorted[2].ID)
	})

	t.Run(`columnMapping == nil`, func(t *testing.T) {
		_, err := mustUseOneColumnSort(t, db, t.Name(), sqlutil.OneColumnSort{}, nil)
		assert.Error(t, err)
	})
}

func mustUseOneColumnSort(t *testing.T, db *storage.DB, testCase string, sort sqlutil.OneColumnSort, columnMapping map[string]string) ([]testRow, error) {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q, err := sqlutil.UseOneColumnSort(q, sort, columnMapping)
	if err != nil {
		return nil, err
	}

	sql, args, err := q.ToSql()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return rs, nil
}
