package sqlutil_test

import (
	"testing"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/ptrutil"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/sqlutil"
	"github.com/stretchr/testify/assert"
)

func TestUseLimitOffsetPagination(t *testing.T) {
	db.Reset()
	resetTestTable(db)

	defaultLimit := 25

	t.Run(`.Limit == 0 and .Offset == 0`, func(t *testing.T) {
		var insert []testRow
		for i := 0; i < 50; i++ {
			insert = append(insert, testRow{IntColumn: ptrutil.Int(i)})
		}

		inserted := mustInsertTestRows(t, db, t.Name(), insert)

		paginated := mustUseLimitOffsetPagination(t, db, t.Name(), sqlutil.LimitOffsetPagination{})

		if !assert.Len(t, paginated, defaultLimit) {
			t.FailNow()
		}

		for i := 0; i < defaultLimit; i++ {
			assert.Equal(t, inserted[i].ID, paginated[i].ID)
		}
	})

	t.Run(`.Limit == 5 and .Offset == 0`, func(t *testing.T) {
		var insert []testRow
		for i := 0; i < 50; i++ {
			insert = append(insert, testRow{IntColumn: ptrutil.Int(i)})
		}

		inserted := mustInsertTestRows(t, db, t.Name(), insert)

		paginated := mustUseLimitOffsetPagination(t, db, t.Name(), sqlutil.LimitOffsetPagination{
			Limit: 5,
		})

		if !assert.Len(t, paginated, 5) {
			t.FailNow()
		}

		for i := 0; i < 5; i++ {
			assert.Equal(t, inserted[i].ID, paginated[i].ID)
		}
	})

	t.Run(`.Limit == 0 and .Offset == 5`, func(t *testing.T) {
		var insert []testRow
		for i := 0; i < 50; i++ {
			insert = append(insert, testRow{IntColumn: ptrutil.Int(i)})
		}

		inserted := mustInsertTestRows(t, db, t.Name(), insert)

		paginated := mustUseLimitOffsetPagination(t, db, t.Name(), sqlutil.LimitOffsetPagination{
			Offset: 5,
		})

		if !assert.Len(t, paginated, defaultLimit) {
			t.FailNow()
		}

		for i := 0; i < defaultLimit; i++ {
			assert.Equal(t, inserted[i+5].ID, paginated[i].ID)
		}
	})
}

func mustUseLimitOffsetPagination(t *testing.T, db *storage.DB, testCase string, pagination sqlutil.LimitOffsetPagination) []testRow {
	q := sqlutil.Select("*").From("sqlutil_test")

	// only include rows that have been created for this testcase
	q = sqlutil.UseStringFilter(q, "test_case", sqlutil.StringFilter{
		Is: &testCase,
	})

	q = sqlutil.UseLimitOffsetPagination(q, pagination)

	sql, args, err := q.ToSql()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var rs []testRow
	err = db.Select(&rs, sql, args...)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return rs
}
